const password = (v: string) =>
  /^(?=.*[a-zA-Z])(?=.*[!@#$%^*+=-])(?=.*[0-9]).{8,16}$/.test(v) ||
  self.$nuxt.$t("msgPassword");

export default password;
