<template>
  <v-container id="page-login" class="login-container fill-height">
    <v-row>
      <v-col cols="4" class="login-logo d-flex align-content-center flex-wrap">
        <v-row height="300">
          <img
            class="login-box-img"
            src="/images/logo_cb.png"
            height="56"
            width="150"
            alt="Cloud-Barista API G/W Admin"
          />
          <div class="login-box-admin">
            API G/W ADMIN
          </div>
          <span class="login-box-admin-span">2020 v1.0</span>
        </v-row>
        <v-row>
          <v-col cols="12" align="end" justify="end">
            <cb-lang-select
              :languages="languages"
              :current-language="currentLanguage"
              @languageChanged="languageChanged"
            />
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="8" class="login-box ml-0 pa-8" height="300">
        <login-form ref="loginForm" v-model="loginModel" @submit="doLogin" />
      </v-col>
    </v-row>
    <cb-loading
      :loading="loading"
      :value="$t('msgLogin')"
      color="red"
      size="150"
    ></cb-loading>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue, Ref } from "nuxt-property-decorator";

import LoginForm from "@/views/components/composite/login-form.vue";
import { appStore } from "@/store";
import { ILogin } from "@/models";

@Component({
  layout: "none",
  components: {
    "login-form": LoginForm
  },
  auth: "guest"
})
export default class LoginPage extends Vue {
  @Ref() loginForm!: LoginForm;

  // ---------------------------------
  // Fields
  // ---------------------------------

  private loading: boolean = false;
  private alert = { style: "error", message: "" };
  private login = { userId: "", password: "" };
  private loginModel: ILogin = { username: "", password: "" };

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get languages() {
    return appStore.languages;
  }

  private get currentLanguage() {
    return appStore.locale;
  }

  // ---------------------------------
  // Methdos
  // ---------------------------------

  private languageChanged(id: string) {
    appStore.setLocale(id);
    this.loginForm.reset();
  }

  private doLogin() {
    this.loading = true;
    this.$auth
      .loginWith("local", { data: this.loginModel })
      .then(_ => {
        // console.log(`LOGIN Response: ${JSON.stringify(res)}`);
      })
      .catch(err => {
        if (err.message.includes("401")) {
          this.$dialog.warning({
            text: this.$t("errLogin") as string,
            title: "로그인"
          });
        }
      })
      .finally(() => {
        this.loading = false;
      });
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------

  created() {
    // if (this.$auth.$state.loggedIn) {
    //   this.$router.replace({ name: "home" });
    // }
  }

  mounted() {}
}
</script>

<style lang="scss" scoped>
.login-container {
  position: relative;
  margin: 0 auto;
  min-height: calc(100vh - 64px - 50px - 1px);
  max-width: 800px;
  //background-color: $login_bg;
}

.login-logo {
  position: relative;
  overflow: hidden;
  text-align: center;
  margin: 0 auto;
  width: 100%;
}

.login-box-img {
  position: relative;
  text-align: center;
  margin: 0 auto;
}

.login-box-admin {
  position: relative;
  text-align: right;
  margin: 0 auto;
  width: 210px;
  font-size: 1.7em !important;
  font-weight: 500;
  margin-top: 5px;
}

.login-box-admin-span {
  position: relative;
  margin: 0 auto;
  text-align: right;
  width: 210px;
  font-size: 13px;
  font-weight: 300;
  margin-top: -2px;
}

.login-box {
  position: relative;
  overflow: hidden;
  text-align: center;
  margin: 0;
  box-sizing: border-box;
  -o-box-sizing: border-box;
  -ms-box-sizing: border-box;
  -moz-box-sizing: border-box;
  -webkit-box-sizing: border-box;
}

.login-card {
  background-color: transparent;
  padding: 5px;
}

.v-card__text {
  padding: 4px !important;
  line-height: 1em;
}

/* dark */
.theme--dark .login-logo {
  background-color: #33519e;
}

.theme--dark .login-box {
  background-color: #30333c;
}

/* light */
.theme--light .login-logo {
  background-color: #33519e;
}

.theme--light .login-box-admin {
  color: #a1b1c1;
}

.theme--light .login-box-admin-span {
  color: #a1b1c1;
}

.theme--light .login-box {
  background-color: #a1b1c1;
}
</style>
