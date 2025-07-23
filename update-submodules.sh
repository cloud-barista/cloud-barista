#!/bin/bash

# README for this script
# This script is used for updating Git submodules to a specific tag and optionally committing those changes.
# It reads a list of submodule URLs from an external file, checks each submodule, updates it to a specified tag,
# and then commits the change if the user agrees.

# Usage:
# 1. Create a file (default is 'submodules.md') and list each submodule URL on a new line.
# 2. Run this script. It will read the submodule URLs from the file.
# 3. For each submodule, the script will:
#    - Check if the submodule directory exists. If not, it will prompt to add the submodule.
#    - Fetch all tags from the remote repository.
#    - Prompt the user to select a tag for checkout (default is the latest tag).
#    - Checkout the specified tag.
#    - Prompt the user to commit the change. If agreed, it will commit the change.
# 4. After processing all submodules, the script ends with a message indicating the completion.

# Important:
# - Ensure you have the necessary permissions to modify the submodules and push changes.
# - Verify the submodule URLs in your file are correct and accessible.

# Debug function to check Git repository status
debug_git_status() {
    local submodule_path=$1
    echo ""
    echo "=== Debug Info for $submodule_path ==="
    echo "Working directory: $(pwd)"
    echo "Remote URL: $(git remote get-url origin 2>/dev/null || echo 'No remote found')"
    echo "Current branch: $(git branch --show-current 2>/dev/null || echo 'Detached HEAD')"
    echo "Latest commit: $(git log -1 --oneline 2>/dev/null || echo 'No commits')"
    
    echo "Checking remote connectivity..."
    if git ls-remote --exit-code --heads origin >/dev/null 2>&1; then
        echo "✓ Remote connectivity: OK"
    else
        echo "✗ Remote connectivity: FAILED"
    fi
    
    echo "Remote branches: $(git branch -r 2>/dev/null | wc -l) found"
    echo "Local tags before fetch: $(git tag -l 2>/dev/null | wc -l) found"
    
    echo "Checking remote tags directly..."
    remote_tags=$(git ls-remote --tags origin 2>/dev/null | wc -l)
    echo "Remote tags available: $remote_tags"
    
    if [ "$remote_tags" -gt 0 ]; then
        echo "Latest remote tags:"
        git ls-remote --tags origin | grep 'refs/tags/v' | tail -5 | sed 's/.*refs\/tags\//  /'
    fi
    
    echo "=== End Debug Info ==="
    echo ""
}


# Function to add a new submodule
add_new_submodule() {
    local submodule_url=$1
    local submodule_path=$(basename $submodule_url .git)

    read -r -p "The submodule $submodule_path does not exist. Do you want to add it? (yes/NO): " add_answer
    if [[ $add_answer == "yes" ]]; then
        git submodule add "$submodule_url" "$submodule_path"
        if [ $? -ne 0 ]; then
            echo "Error: Failed to add submodule $submodule_path"
            return 1
        fi
        echo "Submodule $submodule_path added successfully."
    else
        echo "Skipping addition of $submodule_path."
        return 1
    fi
}

# Function to update a specific submodule to a target tag and optionally commit the change
update_submodule() {
    # Extract the submodule name from the URL
    local submodule_url=$1
    local submodule_path=$(basename $submodule_url .git)

    # Check if the submodule directory exists
    if [[ ! -d "$submodule_path" ]]; then
        # Attempt to add the new submodule
        add_new_submodule "$submodule_url" || return
    fi

    # Enter the submodule directory or exit on failure
    cd "$submodule_path" || { echo "Error: Failed to enter $submodule_path directory."; exit 1; }

    # Check if this is a properly initialized git repository
    if [ ! -d ".git" ]; then
        echo "Warning: $submodule_path is not a proper git repository"
        echo "Attempting to reinitialize submodule..."
        cd ..
        git submodule update --init --recursive "$submodule_path"
        cd "$submodule_path"
    fi

    # Debug git status
    debug_git_status "$submodule_path"

    # Fetch all tags from the remote repository
    echo "Checking and updating remote repository connection..."
    
    # Ensure remote origin is properly set
    if ! git remote get-url origin >/dev/null 2>&1; then
        echo "Error: No remote origin found for $submodule_path"
        git remote -v
        exit 1
    fi
    
    echo "Remote URL: $(git remote get-url origin)"
    
    # Update remote URL to ensure it's accessible
    remote_url=$(git remote get-url origin)
    git remote set-url origin "$remote_url"
    
    echo "Fetching all remote data (this may take a moment)..."
    
    # Try multiple fetch strategies
    if ! git fetch origin --tags --force --prune; then
        echo "Warning: Standard fetch failed, trying alternative methods..."
        
        # Try fetching specific refs
        if ! git fetch origin '+refs/heads/*:refs/remotes/origin/*' '+refs/tags/*:refs/tags/*'; then
            echo "Warning: Alternative fetch also failed, trying basic fetch..."
            git fetch origin
        fi
    fi
    
    # Explicitly fetch all remote references
    echo "Updating all remote references..."
    git remote update origin --prune
    
    # Force update local tag cache
    git tag -d $(git tag -l) 2>/dev/null || true
    git fetch origin --tags --force
    
    # Show remote info for debugging
    echo "Remote repository: $(git remote get-url origin)"
    
    # Get all tags and show count
    echo "Collecting all available tags..."
    all_tags=$(git tag -l | sort -V)
    tag_count=$(echo "$all_tags" | grep -v '^$' | wc -l)
    echo "Local tags after fetch: $tag_count"
    
    # Also check remote tags one more time
    remote_tag_count=$(git ls-remote --tags origin 2>/dev/null | grep 'refs/tags/' | wc -l)
    echo "Remote tags confirmed: $remote_tag_count"
    
    if [ "$tag_count" -eq 0 ] && [ "$remote_tag_count" -gt 0 ]; then
        echo "Warning: Remote has tags but local doesn't. Trying forced fetch..."
        git fetch origin '+refs/tags/*:refs/tags/*' --force
        all_tags=$(git tag -l | sort -V)
        tag_count=$(echo "$all_tags" | grep -v '^$' | wc -l)
        echo "Tags after forced fetch: $tag_count"
    fi
    
    # Determine the latest tag using multiple strategies
    # Strategy 1: Try semantic versioning pattern (v0.0.0)
    latest_tag_v=$(git tag -l | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' | sort -V | tail -1)
    
    # Strategy 2: Try broader pattern (v0.0)  
    latest_tag_broad=$(git tag -l | grep -E '^v[0-9]+\.[0-9]+' | sort -V | tail -1)
    
    # Strategy 3: All tags starting with v
    latest_tag_v_all=$(git tag -l | grep -E '^v' | sort -V | tail -1)
    
    # Strategy 4: All tags
    latest_tag_all=$(git tag -l | sort -V | tail -1)
    
    # Choose the best latest tag
    if [[ -n "$latest_tag_v" ]]; then
        latest_tag="$latest_tag_v"
        echo "Using semantic version pattern: $latest_tag"
    elif [[ -n "$latest_tag_broad" ]]; then
        latest_tag="$latest_tag_broad"
        echo "Using broad version pattern: $latest_tag"
    elif [[ -n "$latest_tag_v_all" ]]; then
        latest_tag="$latest_tag_v_all"
        echo "Using v-prefixed tags: $latest_tag"
    elif [[ -n "$latest_tag_all" ]]; then
        latest_tag="$latest_tag_all"
        echo "Using all tags: $latest_tag"
    else
        # Determine default branch (main vs master)
        default_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || echo "main")
        if ! git show-ref --verify --quiet "refs/remotes/origin/$default_branch"; then
            default_branch="master"
        fi
        latest_tag="$default_branch"
        echo "Warning: No tags found for $submodule_path, will use $default_branch branch"
    fi

    # Sort and display the list of available tags
    echo ""
    echo "[Available tags in $submodule_path]"
    if [[ -n "$all_tags" ]]; then
        echo "$all_tags"
    else
        echo "No tags found"
    fi
    echo ""
    echo "Latest detected tag: $latest_tag"

    # Prompt user to select a tag, prefilling with the latest tag
    echo ""
    read -r -p "Select a tag to checkout for $submodule_path (Just ENTER for $latest_tag): " TARGET_TAG
    TARGET_TAG=${TARGET_TAG:-$latest_tag}

    # Validate the selected tag exists
    if [[ "$TARGET_TAG" != "main" ]] && [[ "$TARGET_TAG" != "master" ]] && ! git show-ref --verify --quiet "refs/tags/$TARGET_TAG" && ! git show-ref --verify --quiet "refs/remotes/origin/$TARGET_TAG"; then
        echo "Error: Tag or branch '$TARGET_TAG' does not exist in $submodule_path"
        echo "Available tags (latest 10):"
        git tag -l | sort -V | tail -10
        echo "Available branches:"
        git branch -r | head -5
        exit 1
    fi

    # Checkout to the tag specified by the user
    echo "Attempting to checkout: $TARGET_TAG"
    
    if [[ "$TARGET_TAG" == "main" ]] || [[ "$TARGET_TAG" == "master" ]]; then
        git checkout "$TARGET_TAG"
        git pull origin "$TARGET_TAG"
    else
        # Try different checkout methods
        if git show-ref --verify --quiet "refs/tags/$TARGET_TAG"; then
            git checkout "tags/$TARGET_TAG"
        elif git show-ref --verify --quiet "refs/remotes/origin/$TARGET_TAG"; then
            git checkout "$TARGET_TAG"
        else
            echo "Error: Could not find tag or branch '$TARGET_TAG'"
            echo "Available references:"
            git tag -l | head -10
            echo "..."
            exit 1
        fi
    fi
    
    if [ $? -ne 0 ]; then
      echo "Error: Failed to checkout $TARGET_TAG in $submodule_path"
      exit 1
    fi

    echo "Successfully checked out $TARGET_TAG in $submodule_path"

    # Ask the user if they want to commit the change
    read -r -p "Do you want to commit the change for $submodule_path? (ENTER/no): " commit_answer

    if [[ $commit_answer != "no" ]]; then
        # Navigate back to the parent directory
        cd ..

        # Stage the changes for the submodule
        git add "$submodule_path"

        # Create a commit
        git commit -m "Update submodule $submodule_path to $TARGET_TAG"
        echo "$submodule_path successfully updated to $TARGET_TAG and committed."
    else
        echo "$submodule_path successfully updated to $TARGET_TAG without committing."
    fi
}

# Read submodules from an external file
submodules_file="submodules.md"

# Check if the file exists
if [[ ! -f $submodules_file ]]; then
    echo "Error: File '$submodules_file' not found."
    exit 1
fi

# Read submodules into an array
readarray -t submodule_urls < "$submodules_file"

# Update each submodule
for submodule_url in "${submodule_urls[@]}"; do
    update_submodule "$submodule_url"
done

echo "All submodules have been successfully updated."
echo "ToDo: git log"
echo "ToDo: git push origin"