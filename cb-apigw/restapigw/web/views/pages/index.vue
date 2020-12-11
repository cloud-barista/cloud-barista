<template>
  <v-container fluid fill-height>
    <v-row>
      <v-col cols="12">
        <v-card class="mx-auto" max-width="70%">
          <v-card-title>
            <v-row dense>
              <v-col class="mx-auto">
                <h2>API Groups</h2>
              </v-col>
            </v-row>
          </v-card-title>
          <v-divider />
          <v-card-text>
            <v-row dense>
              <v-col dense class="text-right">
                <cb-button
                  color="primary"
                  icon="mdi-refresh"
                  class="mr-2"
                  @click="refreshGroup()"
                >
                  Refresh Groups
                </cb-button>
                <cb-button
                  class="mr-2"
                  color="primary"
                  icon="mdi-plus"
                  @click="newGroup()"
                >
                  New Group
                </cb-button>
                <cb-button
                  color="primary"
                  icon="mdi-file-import"
                  @click="$refs.file.click()"
                >
                  Load Group
                </cb-button>
              </v-col>
            </v-row>
            <v-row dense>
              <v-col dense>
                <v-expansion-panels :value="expanded">
                  <v-expansion-panel
                    v-for="(group, index) in apiGroups"
                    :key="index"
                  >
                    <v-expansion-panel-header>
                      <v-row dense>
                        <v-col cols="4" dense>
                          {{ group.name }}
                        </v-col>
                        <v-col dense class="text-right mx-5">
                          <v-icon class="mr-2" @click.stop="viewGroup(group)">
                            mdi-eye
                          </v-icon>
                          <v-icon class="mr-2" @click.stop="exportGroup(group)">
                            mdi-content-save
                          </v-icon>
                          <v-icon
                            color="warning"
                            @click.stop="removeGroup(group)"
                          >
                            mdi-delete
                          </v-icon>
                        </v-col>
                      </v-row>
                    </v-expansion-panel-header>
                    <v-expansion-panel-content>
                      <v-data-table
                        :headers="headers"
                        :items="group.definitions"
                        hide-default-footer
                        sort-by="name"
                        class="elevation-1"
                        fixed-header
                      >
                        <template v-slot:top>
                          <v-toolbar flat dense>
                            <v-toolbar-title>Apis</v-toolbar-title>
                            <v-divider class="mx-4" inset vertical />
                            <v-spacer />
                            <cb-button
                              color="primary"
                              icon="mdi-refresh"
                              class="mr-2"
                              @click="refreshApi(group.name)"
                            >
                              Refresh APIs
                            </cb-button>
                            <cb-button
                              color="primary"
                              icon="mdi-plus"
                              @click="newApi(group.name)"
                            >
                              New API
                            </cb-button>
                          </v-toolbar>
                        </template>
                        <template v-slot:[`item.active`]="{ item }">
                          <v-simple-checkbox
                            v-model="item.active"
                            @click="activeChange(group.name, item)"
                          ></v-simple-checkbox>
                        </template>
                        <template v-slot:[`item.actions`]="{ item }">
                          <v-icon
                            small
                            class="mr-2"
                            @click="viewDefinition(group.name, item)"
                          >
                            mdi-eye
                          </v-icon>
                          <v-icon
                            class="mr-2"
                            small
                            @click.stop="editDefinition(group.name, item)"
                          >
                            mdi-file-edit
                          </v-icon>
                          <v-icon
                            small
                            class="mr-2"
                            color="warning"
                            @click="removeDefinition(group.name, item)"
                          >
                            mdi-delete
                          </v-icon>
                        </template>
                      </v-data-table>
                    </v-expansion-panel-content>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-col>
            </v-row>
          </v-card-text>
          <v-divider />
          <v-card-actions>
            <v-row>
              <v-col class="text-right">
                <cb-button
                  class="mr-2"
                  color="success"
                  icon="mdi-content-save-all"
                  @click="applyChanges"
                >
                  Appy Changes
                </cb-button>
              </v-col>
            </v-row>
          </v-card-actions>
        </v-card>
      </v-col>
      <!-- Dialogs -->
      <yaml-dialog :data="yamlData" @cancel="dialogCancel" @ok="dialogOK" />
      <!-- file -->
      <input
        ref="file"
        type="file"
        style="display:none;"
        @change="importGroup"
      />
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Component, Vue } from "nuxt-property-decorator";
import { YamlDialog } from "@/views/components";

import Util from "@/utils";
import { ApiDefinition, ApiGroup, BackendConfig } from "@/models";
import { apiStore } from "@/store";
import { search } from "@/utils/validation";

@Component({
  auth: true,
  layout: "default",
  name: "index-page",
  components: {
    "yaml-dialog": YamlDialog
  }
})
export default class IndexPage extends Vue {
  // ---------------------------------
  // Fields
  // ---------------------------------

  private searchWord: string = "";
  private headers = [
    {
      text: "Name",
      sortable: false,
      align: "start",
      width: "50%",
      value: "name"
    },
    {
      text: "Endpoint",
      sortable: false,
      align: "start",
      value: "endpoint"
    },
    {
      text: "Active",
      sortable: false,
      align: "center",
      width: "100px",
      value: "active"
    },
    {
      text: "Actions",
      sortable: false,
      align: "center",
      width: "140px",
      value: "actions"
    }
  ];
  private apiGroups: Array<ApiGroup> = [];
  private expanded = [];
  private yamlData = {
    title: "",
    dialog: false,
    action: "",
    content: "",
    readonly: true,
    type: "group",
    gname: ""
  };

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get searchRule() {
    return search;
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private dialogOK(item: any) {
    switch (item.type) {
      case "definition":
        if (item.action === "new") {
          const error = (item.item as ApiGroup).Validate();
          if (error !== "") {
            this.$dialog.error({ title: "New Api Definition", text: error });
          } else {
            apiStore
              .AddApiDefinition(item.item)
              .then(() => {
                this.$dialog.notify.info("API Definition이 추가되었습니다.");
                this.$nextTick(() => {
                  this.apiGroups = apiStore.groups;
                });
              })
              .catch((_: any) => {});
          }
        } else if (item.action === "update") {
          const error = (item.item as ApiGroup).Validate();
          if (error !== "") {
            this.$dialog.error({ title: "Update Api Definition", text: error });
          } else {
            apiStore
              .UpdateApiDefinition(item.item)
              .then(() => {
                this.$dialog.notify.info("API Definition이 변경되었습니다.");
                this.$nextTick(() => {
                  this.apiGroups = apiStore.groups;
                });
              })
              .catch((_: any) => {});
          }
        }
        this.yamlData.dialog = false;
        break;
      case "group":
      default:
        this.yamlData.dialog = false;
        break;
    }
  }

  private dialogCancel(item: any) {
    switch (item.type) {
      case "group":
      case "definition":
      default:
        this.yamlData.dialog = false;
        break;
    }
  }

  private activeChange(gname: string, item: ApiDefinition) {
    const group: ApiGroup = new ApiGroup();
    group.name = gname;
    group.definitions.push(item);
    group.AdjustSendValues();

    apiStore
      .UpdateApiDefinition(group)
      .then(() => {
        this.$dialog.notify.info(
          "API Definition의 활성화 여부가 반영되었습니다."
        );
      })
      .catch((_: any) => {
        item.active = !item.active;
      });
  }

  private refreshGroup() {
    apiStore.GetApiGroups().then(() => {
      this.apiGroups = apiStore.groups;
      if (!this.apiGroups) {
        this.$dialog.notify.info("관리중인 API GROUP 정보가 없습니다.");
      } else {
        this.$dialog.notify.info("API GROUP 정보가 조회되었습니다.");
      }
    });
  }

  private newGroup() {
    this.$dialog
      .prompt({
        text: "Group 명을 입력하십시오.",
        title: "신규 Group 생성"
      })
      .then(data => {
        if (!data) return;

        // validate group
        if (this.apiGroups.filter(g => g.name === data).length >= 1) {
          this.$dialog.warning({
            title: "New Group",
            text: "이미 존재하는 Group 입니다."
          });
        } else {
          const group: ApiGroup = {
            name: data,
            definitions: [] as Array<ApiDefinition>
          } as ApiGroup;
          apiStore.AddApiGroup(group).then(() => {
            this.$dialog.notify.info("API GROUP이 등록되었습니다.");
            this.$nextTick(() => {
              this.apiGroups = apiStore.groups;
            });
          });
        }
      });
  }

  private importGroup(ev: any) {
    const file = ev.target.files[0];
    const reader = new FileReader();
    reader.onload = e => {
      // 파일 내용을 Model 구조로 재 구성
      const group = this.rebuildGroup(
        Util.deserializeYaml(e.target!.result as string)
      );
      group.name = file.name;

      // 동일한 그룹 존재여부 검증
      if (this.apiGroups.filter(g => g.name === group.name).length > 0) {
        this.$dialog.warning({
          text:
            "이미 존재하는 그룹 정보입니다.<br/> 삭제 후 다시 등록하십시오.",
          title: "Load Group"
        });
      } else {
        group.AdjustSendValues();
        const error = group.Validate();
        if (error !== "") {
          this.$dialog.error({ title: "Load API Gorup", text: error });
        } else {
          apiStore
            .AddApiGroup(group)
            .then(() => {
              this.$dialog.notify.info("API GROUP이 등록되었습니다.");
              this.$nextTick(() => {
                this.apiGroups = apiStore.groups;
              });
            })
            .catch(_ => {});
        }
      }
      // Input file 초기화
      (this.$refs.file as HTMLInputElement).value = "";
    };

    reader.readAsText(file);
  }

  private viewGroup(group: ApiGroup) {
    this.yamlData.type = "group";
    this.yamlData.title = "Group Definition (View)";
    this.yamlData.dialog = true;
    this.yamlData.content = Util.serializeYaml(group);
    this.yamlData.readonly = true;
  }

  private exportGroup(group: ApiGroup) {
    const data = Util.serializeYaml(group);
    const blob = new Blob([data], { type: "text/yaml" });
    const e = document.createEvent("MouseEvents");
    const a = document.createElement("a");
    a.download = group.name;
    a.href = window.URL.createObjectURL(blob);
    a.dataset.downloadurl = ["text/yaml", a.download, a.href].join(":");
    e.initMouseEvent(
      "click",
      true,
      false,
      window,
      0,
      0,
      0,
      0,
      0,
      false,
      false,
      false,
      false,
      0,
      null
    );
    a.dispatchEvent(e);
  }

  private removeGroup(group: ApiGroup) {
    this.$dialog
      .confirm({
        title: "Remove Group",
        text: "해당 그룹의 모든 정보가 삭제됩니다.<br/>삭제하시겠습니까?"
      })
      .then(data => {
        if (data) {
          apiStore
            .RemoveApiGroup(group.name)
            .then(() => {
              this.$dialog.notify.info("API GROUP이 삭제되었습니다.");
              this.$nextTick(() => {
                this.apiGroups = apiStore.groups;
              });
            })
            .catch(_ => {});
        }
      });
  }

  private refreshApi(gname: string) {
    apiStore
      .GetApiGroup(gname)
      .then(() => {
        this.$dialog.notify.info("API Definition이 조회되었습니다.");
        this.$nextTick(() => {
          this.apiGroups = apiStore.groups;
        });
      })
      .catch(_ => {});
  }

  private newApi(gname: string) {
    const def = new ApiDefinition();

    this.yamlData.type = "definition";
    this.yamlData.title = "API Definition (New)";
    this.yamlData.action = "new";
    this.yamlData.dialog = true;
    this.yamlData.content = Util.serializeYaml(def);
    this.yamlData.readonly = false;
    this.yamlData.gname = gname;
  }

  private viewDefinition(gname: string, item: ApiDefinition) {
    this.yamlData.type = "definition";
    this.yamlData.title = "API Definition (View)";
    this.yamlData.dialog = true;
    this.yamlData.content = Util.serializeYaml(item);
    this.yamlData.readonly = true;
    this.yamlData.gname = gname;
  }

  private editDefinition(gname: string, item: ApiDefinition) {
    this.yamlData.type = "definition";
    this.yamlData.title = "API Definition (Edit)";
    this.yamlData.action = "update";
    this.yamlData.dialog = true;
    this.yamlData.content = Util.serializeYaml(item);
    this.yamlData.readonly = false;
    this.yamlData.gname = gname;
  }

  private removeDefinition(gname: string, item: ApiDefinition) {
    this.$dialog
      .confirm({
        title: "Remove Definition",
        text: "API Definition 정보가 삭제됩니다.<br/>삭제하시겠습니까?"
      })
      .then(data => {
        if (data) {
          const group = new ApiGroup();
          group.name = gname;
          group.definitions = [item];

          apiStore
            .RemoveApiDefinition(group)
            .then(() => {
              this.$dialog.notify.info("API Definition이 삭제되었습니다.");
              this.$nextTick(() => {
                this.apiGroups = apiStore.groups;
              });
            })
            .catch(_ => {});
        }
      });
  }

  private applyChanges() {
    this.$dialog
      .confirm({
        title: "Apply Changes",
        text:
          "모든 변경된 정보가 전체 시스템에 반영됩니다.<br/>적용하시겠습니까?"
      })
      .then(data => {
        if (data) {
          apiStore
            .ApplyChanges()
            .then(() => {
              this.$dialog.notify.info(
                "API Groups 변경 내용이 적용되었습니다."
              );
              this.$nextTick(() => {
                this.apiGroups = apiStore.groups;
              });
            })
            .catch(_ => {});
        }
      });
  }

  private rebuildGroup(obj: any): ApiGroup {
    // 파일 내용을 Model 구조로 재 구성
    const group = Object.assign(new ApiGroup(), obj);
    group.definitions = group.definitions.map((d: any) => {
      const def = Object.assign(new ApiDefinition(), d) as ApiDefinition;
      def.backend = def.backend.map((b: any) =>
        Object.assign(new BackendConfig(), b)
      );
      return def;
    });
    return group;
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------

  created() {
    this.$nextTick(() => this.refreshGroup());
  }
}
</script>

<style lang="scss" scoped></style>
