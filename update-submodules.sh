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

    # Fetch all tags from the remote repository
    git fetch --tags

    # Determine the latest tag
    latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

    # Sort and display the list of available tags
    echo "[Available tags in $submodule_path]"
    git tag -l | sort -V

    # Prompt user to select a tag, prefilling with the latest tag
    read -r -p "Select a tag to checkout for $submodule_path (Just ENTER for $latest_tag): " TARGET_TAG
    TARGET_TAG=${TARGET_TAG:-$latest_tag}

    # Checkout to the tag specified by the user
    git checkout "tags/$TARGET_TAG"
    if [ $? -ne 0 ]; then
      echo "Error: Failed to checkout tag $TARGET_TAG in $submodule_path"
      exit 1
    fi

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