// Vuetify Colors 변경
// import colors from "vuetify/lib/util/colors";

const apigw = {
  title: "Cloud-Barista REST API G/W Admin", // 브라우저에 표시될 타이틀
  host: "0.0.0.0", // 개발 검증용 (로컬 테스트 구동에만 사용 - Node 기반)
  port: 4444, // 개발 검증용 (로컬 테스트 구동에만 사용 - Node 기반)
  api:
    process.env.NODE_ENV === "production"
      ? "" // Admin Web을 분리해서 서비스할 경우는 실제 API G/W Admin API URL을 지정해야 한다. (현재는 API G/W에서 Admin Web/API 동일하게 제공)
      : "http://localhost:8001", // 개발 검증용 (로컬 테스트 구동)
  path: "/"
};

const isDev = process.env.NODE_ENV !== "production";

export default {
  compress: true,
  /**
   * Server Configuration (Dev)
   */
  server: {
    port: apigw.port,
    host: apigw.host
  },
  /**
   * Application mode
   */
  ssr: false,
  /**
   * Environment variables
   */
  env: {},
  /*
   ** Headers of the page
   */
  head: {
    title: apigw.title || "REST API G/W Admin",
    meta: [
      { charset: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      {
        hid: "description",
        name: "description",
        content: process.env.npm_package_description || ""
      }
    ],
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }],
    // Scripts for application in HEAD
    script: [
      { src: "/scripts/yaml.js" } // Yamljs
    ]
  },
  /**
   * Folders
   */
  dir: {
    assets: "utils/assets",
    middleware: "utils/middleware",
    plugins: "utils/plugins",
    static: "utils/static",
    components: "views/components",
    layouts: "views/layouts",
    pages: "views/pages"
  },
  /**
   ** SpinKit for loading Indicator
   */
  loadingIndicator: {
    name: "circle",
    color: "#3B8070",
    background: "white"
  },
  /*
   ** Customize the progress-bar color
   */
  loading: { color: "#3adced" },
  /*
   ** Global CSS
   */
  css: [
    "@/utils/assets/scss/index.scss",
    // Vuetify - Roboto fonts
    "@/utils/assets/fonts/css/roboto.css",
    // Vueitfy - Material Design Icon
    "@/utils/assets/fonts/css/materialdesignicons.min.css"
  ],
  /*
   ** Plugins to load before mounting the App
   */
  plugins: ["~/utils/plugins/i18n.js", "~/utils/plugins/axios-accessor.ts"],
  /*
   ** Nuxt.js modules
   */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    "@nuxtjs/axios",
    // Doc: https://github.com/yariksav/vuetify-dialog
    "vuetify-dialog/nuxt",
    // Doc: https://auth.nuxtjs.org
    "@nuxtjs/auth"
  ],
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: [
    "@nuxt/typescript-build",
    "@nuxtjs/vuetify",
    "@/_modules/module"
  ],
  /*
   ** Build configuration
   */
  build: {
    // 빌드 속도를 올리기위해 아래 3가지 옵션 추가했습니다. 3개다 아직 Experimental한 기능이므로 문제 발생시 주석처리 해주세요.
    parallel: true,
    cache: true,
    devtools: process.env.NODE_ENV === "development",
    analyze: process.env.NODE_ENV === "development",
    // analyze: {
    //   analyzerMode: 'server',
    //   analyzerHost: '0.0.0.0',
    //   analyzerPort: '8888',
    //   openAnalyzer: true
    // },

    // CSS 연관으로 설정하면 Production 모드에서 제대로 동작하지 않는 문제 발생
    // cssSourceMap: true,
    // extractCSS: process.env.NODE_ENV !== 'development',

    transpile: ["vuetify/lib"],
    /*
     ** You can extend webpack config here
     */
    extend(config: any, ctx: any) {
      // Extend only webpack config for client-bundle
      // Run ESLint on save
      if (ctx.isDev && ctx.isClient) {
        config.devtool =
          process.env.NODE_ENV === "development" ? "#source-map" : "";
        config.module.rules.push({
          enforce: "pre",
          test: /\.(js|ts|vue)$/,
          loader: "eslint-loader",
          exclude: /(node_modules)/
        });
      }
    }
  },
  /*
   ** Typescript runtime lint
   */
  typescript: {
    typeCheck: true,
    loaderOptions: {
      compileOptions: {
        target: "esnext",
        module: "esnext"
      }
    },
    ignoreNotFoundWarings: true
  },
  /*
   ** Nuxt Hooks
   */
  hooks: {},

  /**
   * 아래부터는 각 Module에 대한 설정 추가
   */

  /*
   ** Axios module configuration
   ** See https://axios.nuxtjs.org/options
   */
  axios: {
    debug: isDev,
    timeout: 3000,
    retry: { retries: 0 }, // Axios Intercepts 에서 공통처리하는 메시지 박스에서 무한루프 오류 발생하므로 사용 금지.
    baseURL: (process.env.API_BASE_URL || apigw.api) + apigw.path
    // credentials: true  // cookie 사용 인증인 경우만 처리 (CORS 관련 오류 발생하며 wildcard 사용 불가 오류)
  },
  /*
   ** vuetify module configuration
   ** https://github.com/nuxt-community/vuetify-module
   */
  vuetify: {
    defaultAssets: false,
    customVariables: ["~/utils/assets/scss/variables.scss"],
    icons: {
      iconfont: "mdi"
    },
    theme: {
      dark: true,
      options: {
        customProperties: true
      },
      themes: {
        dark: {
          primary: "#1976D2",
          accent: "#FF4081",
          secondary: "#ffe18d",
          success: "#4CAF50",
          info: "#2196F3",
          warning: "#FB8C00",
          error: "#FF5252",
          background: "#363636"
        },
        light: {
          primary: "#1976D2",
          accent: "#e91e63",
          secondary: "#30b1dc",
          success: "#4CAF50",
          info: "#2196F3",
          warning: "#FB8C00",
          error: "#FF5252",
          background: "#F6F6F6"
        }
      }
    }
  },
  /**
   * Nuxt router configuration
   */
  router: {
    middleware: ["i18n", "auth"]
  },
  /**
   * Nuxt auth configuration
   */
  auth: {
    cookie: false,
    // localStorage: false,
    token: {
      prefix: "token."
    },
    strategies: {
      local: {
        endpoints: {
          login: {
            url: "/auth/login",
            method: "post",
            propertyName: "access_token"
          },
          logout: false,
          user: false
        },
        tokenRequired: true,
        tokenType: "Bearer"
      }
    },
    redirect: {
      login: "/auth/login",
      logout: "/auth/login",
      home: "/"
    }
  },
  /**
   * Vuetify-Dialog configuration
   */
  vuetifyDialog: {
    property: "$dialog",
    confirm: {
      actions: {
        false: {
          text: "취소",
          color: "accent"
        },
        true: {
          text: "확인",
          color: "primary"
        }
      },
      icon: "mdi-help-circle"
    },
    prompt: {
      actions: {
        false: {
          text: "취소",
          color: "accent"
        },
        true: {
          text: "확인",
          color: "primary"
        }
      },
      icon: "mdi-bulletin-board"
    },
    info: {
      actions: {
        true: {
          text: "확인",
          color: "primary"
        }
      },
      icon: "mdi-information"
    },
    error: {
      actions: {
        true: {
          text: "확인",
          color: "primary"
        }
      },
      icon: "mdi-alert-circle"
    },
    warning: {
      actions: {
        true: {
          text: "확인",
          color: "primary"
        }
      },
      icon: "mdi-alert"
    }
  }
};
