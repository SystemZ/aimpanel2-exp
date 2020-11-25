<template>
  <v-card>
    <v-dialog
      v-model="removeFileDialog"
      persistent
      max-width="250px"
    >
      <v-card>
        <v-card-title class="headline">
          Delete file
        </v-card-title>
        <v-card-text>
          <p>Are you sure you want to delete this file?</p>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary darken-1"
            text
            @click="removeFileDialog = false; fileToRemove = {}"
          >
            Cancel
          </v-btn>
          <v-btn
            color="red darken-1"
            text
            @click="remove(fileToRemove)"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card-title>Files</v-card-title>

    <v-container v-if="downloading">
      <v-progress-linear

        v-model="downloadProgress"
        color="amber"
        height="25"
      >
        <strong>{{ Math.ceil(downloadProgress) }}%</strong>
      </v-progress-linear>
    </v-container>

    <v-container v-if="selectedFiles.length > 0">
      <v-btn color="primary" @click="selectedFiles = []">Unselect all</v-btn>
    </v-container>


    <v-list subheader two-line>
      <v-subheader inset>Folders</v-subheader>

      <v-list-item
        @click="goToParentDirectory()"
        v-if="files.selected.info && files.selected.info.name !== files.root.info.name"
      >
        <v-list-item-avatar>
          <v-icon class=" white--text">
            {{ mdiArrowLeft }}
          </v-icon>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Go back</v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-list-item
        :key="item.info.name"
        v-for="item in files.selected.children"
        v-if="item.info.is_dir"
      >
        <v-list-item-avatar @click="files.selected = item">
          <v-icon>
            {{ mdiFolder }}
          </v-icon>
        </v-list-item-avatar>

        <v-list-item-content @click="files.selected = item">
          <v-list-item-title v-text="item.info.name"></v-list-item-title>
          <v-list-item-subtitle
            v-text="prettySize(item.info.size)"></v-list-item-subtitle>
        </v-list-item-content>

        <v-list-item-action>
          <v-btn icon>
            <v-icon>{{ mdiInformation }}</v-icon>
          </v-btn>
        </v-list-item-action>
        <v-list-item-action>
          <v-btn icon @click="fileToRemove = item; removeFileDialog = true"
                 color="red">
            <v-icon>{{ mdiDelete }}</v-icon>
          </v-btn>
        </v-list-item-action>
      </v-list-item>

      <v-divider inset></v-divider>

      <v-subheader inset>Files</v-subheader>

      <v-list-item-group
        v-model="selectedFiles"
        multiple
      >
        <v-list-item
          :key="item.info.name"
          :value="item"
          v-for="item in files.selected.children"
          v-if="!item.info.is_dir"
        >
          <template v-slot:default="{ active }">
            <v-checkbox
              :input-value="active"
            ></v-checkbox>

            <v-list-item-avatar>
              <v-icon>{{ mdiFile }}</v-icon>
            </v-list-item-avatar>

            <v-list-item-content>
              <v-list-item-title v-text="item.info.name"></v-list-item-title>
              <v-list-item-subtitle
                v-text="prettySize(item.info.size)"></v-list-item-subtitle>
            </v-list-item-content>

            <v-list-item-action>
              <v-btn icon>
                <v-icon>{{ mdiInformation }}</v-icon>
              </v-btn>
            </v-list-item-action>

            <v-list-item-action>
              <v-btn icon @click="download(item)"
                     color="blue">
                <v-icon>{{ mdiDownload }}</v-icon>
              </v-btn>
            </v-list-item-action>

            <v-list-item-action>
              <v-btn icon @click="fileToRemove = item; removeFileDialog = true"
                     color="red">
                <v-icon>{{ mdiDelete }}</v-icon>
              </v-btn>
            </v-list-item-action>
          </template>
        </v-list-item>
      </v-list-item-group>
    </v-list>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Node } from '@/types/files';
import {
  mdiArrowLeft,
  mdiFile,
  mdiFolder,
  mdiInformation,
  mdiDelete, mdiDownload
} from '@mdi/js';

interface FileRow {
  icon: string,
  iconClass: string,
  title: string,
  subtitle: string
}

@Component
export default class GsFileManager extends Vue {
  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  serverId !: String;

  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  hostId !: String;

  @Prop({
    type: String, required: true, default: () => {
      return '';
    }
  })
  fileServerAddress !: String;

  serverUrl = '';
  stream = '' as any;
  files = {
    root: {} as Node,
    selected: {} as Node
  };

  removeFileDialog = false;
  fileToRemove = {} as Node

  selectedFiles = [];

  downloading = false;
  downloadProgress = 0;

  //icons
  mdiInformation = mdiInformation;
  mdiFile = mdiFile;
  mdiFolder = mdiFolder;
  mdiArrowLeft = mdiArrowLeft;
  mdiDelete = mdiDelete;
  mdiDownload = mdiDownload;

  mounted() {
    this.serverUrl = '/v1/host/' + this.hostId + '/server/' + this.serverId;
    this.getFiles();
  }

  getFiles() {
    this.$http.get(this.serverUrl + '/file/list').then(res => {
      this.files.root = res.data;
      this.files.selected = this.files.root;
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    });
  }

  // https://stackoverflow.com/a/18650828/1351857
  prettySize(bytes: number): string {
    if (bytes === 0) {
      return '0 B';
    }

    const decimals = 2;
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

    const i = Math.floor(Math.log(bytes) / Math.log(k));

    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  }

  goToParentDirectory() {
    if (this.files.selected.parent_name === this.files.root.info.name) {
      this.files.selected = this.files.root;
    } else {
      let node = this.getParent(this.files.root, this.files.selected.parent_name);
      this.files.selected = node as Node;
    }
  }

  getParent(root: Node, name: string) {
    let node = null;
    if (name === '') {
      return root;
    }

    root.children.some(n => {
      if (n.info.name === name) {
        return node = n;
      }

      if (n.children) {
        return node = this.getParent(n, name);
      }
    });
    return node;
  }

  remove(item: Node) {
    this.$http.delete(this.serverUrl + '/file', {
      data: {
        path: item.path
      }
    }).then(res => {
      setTimeout(() => {
        this.getFiles()
      }, 2000)
    }).catch(e => {
      this.$auth.checkResponse(e.response.status);
    }).finally(() => {
      this.fileToRemove = {} as Node
      this.removeFileDialog = false;
    })
  }

  download(item: Node) {
    this.downloading = true;
    this.$http.post(this.fileServerAddress + '/download', {
      path: item.path
    }, {
      responseType: 'blob',
      onDownloadProgress: (ev) => {
        this.downloadProgress = Math.round((ev.loaded * 100) / ev.total)
        console.log(this.downloadProgress)
      }
    }).then(res => {
      console.log(res)
      const [, filename] = res.headers['content-disposition'].split('filename=');

      let blob = new Blob([res.data], {type: res.headers['content-type']});
      let link = document.createElement('a');
      link.href = window.URL.createObjectURL(blob);
      link.download = filename;
      link.click()
    }).finally(() => {
      this.downloading = false;
    })
  }
}
</script>
