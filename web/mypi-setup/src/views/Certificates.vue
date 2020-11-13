<template>
  <span>
    <v-container fluid>
      <div v-if="!haveData">
        <div class="text-center">
          <v-progress-circular indeterminate />
        </div>
      </div>
      <div v-else-if="certificates.length > 0">
        <v-list>
          <v-subheader>Certificates</v-subheader>
          <v-list-item-group v-model="certificates" color="primary">
            <v-list-item v-for="(item, i) in certificates" :key="i">
              <v-list-item-icon>
                <v-icon v-text="item.icon"></v-icon>
              </v-list-item-icon>
              <v-list-item-content>
                <v-list-item-title v-text="item.text"></v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list-item-group>
        </v-list>
      </div>
      <div v-else>
        <v-row align="center" justify="space-around">
          No Certificates available. You may now import an existing PKI, or
          create a new one.
        </v-row>
        <div class="text-center">
          <v-btn>Import</v-btn>
          <v-btn>Create</v-btn>
        </div>
      </div>
    </v-container>
  </span>
</template>

<script>
import axios from "axios";

export default {
  name: "certificates",
  components: {},
  data: function() {
    return { haveData: false, certificates: [] };
  },
  mounted() {
    axios({ method: "GET", url: "/api/certificates" }).then(
      (result) => {
        var certificates = result.data;
        for (var i = 0; i < certificates.length; i++) {
          certificates[i].uri = "/users/" + certificates[i].text;
        }
        this.certificates = certificates;
        this.haveData = true;
      },
      (error) => {
        if (error != null) {
          error = null;
        }
        this.certificates = [
          { icon: "mdi-certificate", text: "MYPI-Root-CA" },
          { icon: "mdi-certificate", text: "MYPI-Server-CA" },
          { icon: "mdi-certificate", text: "MYPI-Client-CA" },
          { icon: "mdi-certificate-outline", text: "localhost" },
        ];
        //  console.error(error);
      }
    );
  },
};
</script>
