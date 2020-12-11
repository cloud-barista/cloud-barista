const search = (v: string) =>
  v === "" || v.length < 2 ? self.$nuxt.$t("msgSearch") : true;

export default search;
