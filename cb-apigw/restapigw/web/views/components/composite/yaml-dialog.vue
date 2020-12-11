<template>
  <v-dialog :value="show" persistent width="50%">
    <v-card>
      <v-card-title class="headline">
        {{ data.title }}
      </v-card-title>
      <v-card-text>
        <!-- <v-textarea
          ref="yamleditor"
          v-model="content"
          :readonly="readonly"
          full-width
          filled
          no-resize
          height="600px"
          hide-details
        ></v-textarea> -->
        <cb-lined-textarea
          ref="editor"
          v-model="content"
          :disabled="readonly"
          :error-lines="errorLines"
          nowrap
          :styles="{ width: '100%', height: '600px', resize: 'none' }"
        />
      </v-card-text>
      <v-card-actions>
        <v-btn v-if="!readonly" color="primary" @click="complete">
          적용
        </v-btn>
        <v-btn text @click="close">닫기</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { Component, Vue, Prop } from "nuxt-property-decorator";

import Util from "@/utils";
import {
  ApiGroup,
  deserializeGroupFromJSON,
  deserializeDefinitionFromJSON
} from "@/models";

@Component({
  name: "yaml-dialog",
  components: {}
})
export default class YamlDialog extends Vue {
  @Prop({ type: Object }) data!: any;

  // ---------------------------------
  // Fields
  // ---------------------------------

  private errorLines: Array<Number> = [];

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get show() {
    return this.data.dialog;
  }

  private set show(val: boolean) {
    this.data.dialog = val;
  }

  private get content() {
    return this.data.content;
  }

  private set content(val: any) {
    this.data.content = val;
  }

  private get readonly() {
    return this.data.readonly;
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private close() {
    this.$emit("cancel", {
      type: this.data.gname === "" ? "group" : "definition",
      item: undefined
    });
    this.data.content = undefined;
  }

  private showError(err: any) {
    this.errorLines = [err.parsedLine];

    this.$dialog.error({
      title: "YAML Parsing Error",
      text: `${err.message} at ${err.parsedLine} line. [${err.snippet}] `
    });
  }

  private complete() {
    let group: ApiGroup;

    if (this.data.type === "group") {
      try {
        group = deserializeGroupFromJSON(
          Util.deserializeYaml(this.data.content)
        );
      } catch (err) {
        return this.showError(err);
      }
    } else {
      group = new ApiGroup();
      group.name = this.data.gname;
      try {
        group.definitions = [
          deserializeDefinitionFromJSON(Util.deserializeYaml(this.data.content))
        ];
      } catch (err) {
        return this.showError(err);
      }
    }

    this.$emit("ok", {
      type: this.data.type,
      action: this.data.action,
      item: group
    });
    this.data.content = undefined;
    this.errorLines = [];
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------
}
</script>

<style lang="scss" scoped></style>
