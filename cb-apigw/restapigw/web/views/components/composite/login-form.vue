<template>
  <v-card class="login-card">
    <v-card-text>
      <v-form ref="form" @keyup.enter="doLogin">
        <cb-text
          ref="username"
          v-model="value.username"
          outlined
          prependinnericon="mdi-account"
          :rules="[requiredRule]"
          :label="$t('username')"
          autocomplete="$t('username')"
          autofocus
          type="text"
        />
        <cb-password
          ref="password"
          v-model="value.password"
          prependinnericon="mdi-lock"
          :rules="[requiredRule, passwordRule]"
          outlined
          show-view
          required
          autocomplete="current-password"
          @keyup.enter="doLogin"
        />
      </v-form>
    </v-card-text>
    <v-card-actions class="pb-7">
      <v-spacer />
      <cb-button
        block
        color="primary"
        icon="mdi-login"
        :text="$t('login')"
        @click="doLogin"
      />
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import { Component, Vue, Prop, Ref } from "nuxt-property-decorator";

import { required, password } from "@/utils/validation";
import { ILogin } from "@/models";

@Component({
  name: "login-form",
  components: {}
})
export default class LoginForm extends Vue {
  @Prop({ type: Object }) value!: ILogin;
  @Ref() form!: HTMLFormElement;

  // ---------------------------------
  // Fields
  // ---------------------------------

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get requiredRule() {
    return required;
  }

  private get passwordRule() {
    return password;
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private setFocus() {
    (this.$refs.username as any).$el.querySelector("input").focus();
  }

  public reset() {
    this.form.reset();
    this.setFocus();
  }

  private doLogin() {
    if (this.form.validate()) {
      this.$emit("submit");
    }
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------

  mounted() {
    this.$nextTick(() => {
      this.setFocus();
    });
  }
}
</script>

<style lang="scss" scoped></style>
