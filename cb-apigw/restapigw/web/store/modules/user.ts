// =================================================================
//  Defines the module for User
// =================================================================

import {
  createModule /*, mutation, action, Module */
} from "vuex-class-component";

const VuexModule = createModule({
  namespaced: "user",
  strict: false,
  target: "nuxt"
});

export default class UserModule extends VuexModule {
  // ------------------------------------------------
  // Fields
  // ------------------------------------------------
  // private sidebar = true
  // ------------------------------------------------
  // Mutations
  // ------------------------------------------------
  // @mutation
  // private toggle_sidebar(): void {
  //    this.sidebar = !this.sidebar
  // }
  // ------------------------------------------------
  // Actions
  // ------------------------------------------------
  // @action
  // public ToggleTheme(): any {
  //   this.toggle_theme()
  // }
}
