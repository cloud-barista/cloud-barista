// =================================================================
//  Defines the module for Application
// =================================================================

import { createModule, mutation, action } from "vuex-class-component";

// Crate Vuex module
const VuexModule = createModule({
  namespaced: "app",
  strict: false,
  target: "nuxt"
});

// Define Module
export default class AppModule extends VuexModule {
  // ------------------------------------------------
  // Fields
  // ------------------------------------------------

  public locale = "ko";
  public languages: Array<any> = [];

  // ------------------------------------------------
  // Mutations
  // ------------------------------------------------

  @mutation
  private SET_LOCALE(val: string) {
    this.locale = val;
    self.$nuxt.$i18n.locale = val;
  }

  @mutation
  private SET_LANGUAGES(val: Array<any>) {
    this.languages = val;
  }

  // ------------------------------------------------
  // Actions
  // ------------------------------------------------

  @action
  public async setLocale(val: string) {
    this.SET_LOCALE(val);
  }

  @action
  public async setLanguages(val: Array<any>) {
    this.SET_LANGUAGES(val);
  }
}
