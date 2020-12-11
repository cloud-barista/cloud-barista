const { resolve, join } = require("path");
const globby = require("globby");
const chokidar = require("chokidar");

const libRoot = resolve(__dirname, "..");

module.exports = function(moduleOptions) {
  // console.log(`[CB Modules] Module started... LibRoot: ${libRoot}`);

  const { options } = this.nuxt;

  /**
   * Module Options
   */

  // Merge options and apply defaults
  const { dir, suffixes, extensions, ignore, ignoreNameDetection } = {
    // fixes components directory to module's directory
    // dir: options.dir.components || join(libRoot, 'components'),
    dir: join(libRoot, "core/components"),
    suffixes: [""],
    extensions: ["vue", "js", "ts"],
    ignore: [],
    ignoreNameDetection: false,
    ...moduleOptions,
    ...this.options["cb-components"],
    ...this.options.cbComponents
  };

  // Compute suffix with extensions array
  const suffixesWithExtensions = getSuffixesWithExtensions(
    suffixes,
    extensions
  );

  // Remove extension
  const removeExtension = name =>
    name.replace(new RegExp(`\\.(${suffixesWithExtensions.join("|")})$`), "");

  // Resolve components dirs
  const dirs = Array.isArray(dir) ? dir : [dir];
  // if use src directory from nuxt.config.ts
  // const componentsDir = dirs.map(d => resolve(options.srcDir, d))
  const componentsDir = dirs;

  // Patterns
  const patterns = dirs.map(d =>
    join(d, `/**/*.(${suffixesWithExtensions.join("|")})`).replace(/\\/g, "/")
  );

  /**
   * Plugin Options
   */

  // CB-Components plugin options
  const componentOptions = {
    components: [],
    ignoreNameDetection
  };

  // Scans global components and updates context
  const scanGlobalComponents = async () => {
    const fileNames = await globby(patterns, {
      cwd: libRoot,
      ignore,
      absolute: false,
      objectMode: true,
      onlyFiles: true
    });

    const globalComponents = fileNames.map(({ name, path }) => {
      return {
        name: toPascal(removeExtension(name)),
        path: `${path}`
      };
    });

    const changesDetected = !deepEqual(
      globalComponents,
      componentOptions.components
    );
    componentOptions.components = globalComponents;
    return changesDetected;
  };

  /** Nuxt builder's hook */

  // Hook on builder
  this.nuxt.hook("build:before", async builder => {
    // console.log("[CB Components] Nuxt building hook");

    // Scan components once
    await scanGlobalComponents();

    // Watch components directory for dev mode
    if (this.options.dev) {
      const watcher = chokidar.watch(componentsDir, options.watchers.chokidar);
      watcher.on("all", async eventName => {
        if (!["add", "unlink"].includes(eventName)) {
          return;
        }
        const changesDetected = await scanGlobalComponents();
        if (changesDetected) {
          builder.generateRoutesAndFiles();
        }
      });

      // Close watcher on nuxt close
      this.nuxt.hook("close", () => {
        watcher.close();
      });
    }
  });

  /**
   * Plugins
   */
  // // Add common utility functions plugin
  // this.addPlugin({
  //   src: resolve(__dirname, 'cb-utilities.js'),
  //   filename: 'cb-utilities.js',
  //   options: undefined
  // })
  // // Add event-bus plugin
  // this.addPlugin({
  //   src: resolve(__dirname, 'cb-event-bus.js'),
  //   filename: 'cb-event-bus.js',
  //   options: undefined
  // })

  // Add global-components plugin
  this.addPlugin({
    src: resolve(__dirname, "cb-components.js"),
    fileName: "cb-components.js",
    options: componentOptions
  });
};

function getSuffixesWithExtensions(suffixes, extensions) {
  if (suffixes.length < 1) {
    return extensions;
  }

  return suffixes.reduce((acc, suffix) => {
    const suffixWithExtension = extensions.map(
      extension => `${suffix}${suffix ? "." : ""}${extension}`
    );
    return acc.concat(suffixWithExtension);
  }, []);
}

function deepEqual(a, b) {
  return JSON.stringify(a) === JSON.stringify(b);
}

function toPascal(name) {
  const camelCase = name.replace(/([-_]\w)/g, g => g[1].toUpperCase());
  return camelCase[0].toUpperCase() + camelCase.substr(1);
}
