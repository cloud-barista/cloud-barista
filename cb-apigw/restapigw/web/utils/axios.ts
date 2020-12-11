/**
 * Initialize the Axios used in application
 */

import { NuxtAxiosInstance } from "@nuxtjs/axios";

let $axios: NuxtAxiosInstance;

function checkAuth(res: any): Boolean {
  const code = parseInt(res && res.status);
  if (code === 401 || code === 403) {
    self.$nuxt.$auth.logout();
    return false;
  }
  return true;
}
function showError(message: string) {
  self.$nuxt.$dialog.error({ title: "Error", text: message });
}

export function initializeAxios(axiosInstance: NuxtAxiosInstance) {
  // axios request intercepter
  axiosInstance.interceptors.request.use(
    config => {
      // Adds the header to every request, you can add custom headers here.
      // console.log(`axios request = ${JSON.stringify(config)}`)
      return config;
    },
    error => {
      if (checkAuth(error.response)) {
        showError(error.message);
        Promise.reject(error);
      } else {
        return Promise.reject(error);
      }
    }
  );

  // Axios response intercepor
  axiosInstance.interceptors.response.use(
    response => {
      const res = response.data || response;
      if (res && res.code !== 0) {
        showError(res.message);
        return Promise.reject(res);
      }
      return res;
    },
    error => {
      if (checkAuth(error.response)) {
        showError(error.message);
        return Promise.reject(error.message);
      } else {
        return Promise.reject(error);
      }
    }
  );

  $axios = axiosInstance;
}

export { $axios };
