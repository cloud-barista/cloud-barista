import Vue from "vue";
import Vuex from "vuex";
import { createProxy, extractVuexModule } from "vuex-class-component";

// Import store modules
import AppModule from "@/store/modules/app";
import UserModule from "@/store/modules/user";
import ApiModule from "@/store/modules/api";

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    // Add modules to store
    ...extractVuexModule(AppModule),
    ...extractVuexModule(UserModule),
    ...extractVuexModule(ApiModule)
  }
});

// export store from modulees
export const appStore: AppModule = createProxy(store, AppModule);
export const userStore: UserModule = createProxy(store, UserModule);
export const apiStore: ApiModule = createProxy(store, ApiModule);
