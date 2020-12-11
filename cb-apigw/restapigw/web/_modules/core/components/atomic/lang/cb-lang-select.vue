<template>
  <v-menu offset-y>
    <template v-slot:activator="{ on: menu }">
      <v-tooltip bottom>
        <template v-slot:activator="{ on: tooltip }">
          <v-btn icon v-on="{ ...tooltip, ...menu }">
            <img v-if="currentLanguageIcon" :src="currentLanguageIcon" />
            <v-icon v-else>mdi-earth</v-icon>
          </v-btn>
        </template>
        <span>Select Language</span>
      </v-tooltip>
    </template>
    <v-list>
      <v-list-item
        v-for="language in languages"
        :key="language.id"
        @click="changeLanguage(language.id)"
      >
        <v-list-item-avatar tile size="24">
          <v-img :src="language.flagSrc"></v-img>
        </v-list-item-avatar>
        <v-list-item-title>{{ language.title }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

@Component({
  name: "",
  components: {}
})
export default class CBLangSelect extends Vue {
  @Prop({ type: Array, default: () => ["ko", "en"] }) languages!: Array<any>;
  @Prop({ type: String, default: "" }) currentLanguage?: string;

  // ---------------------------------
  // Fields
  // ---------------------------------
  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get currentLanguageIcon() {
    if (!this.currentLanguage) return null;
    return this.languages.filter(x => x.id === this.currentLanguage)[0].flagSrc;
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private changeLanguage(id: string) {
    // this.$router.push(this.switchLocalePath(id));
    this.$emit("languageChanged", id);
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped></style>
