const required = (value: any) => {
  const data = value === 0 ? String(value) : value;
  return !!data || self.$nuxt.$t("msgRequired", self.$nuxt.$i18n.locale);
};

export default required;
