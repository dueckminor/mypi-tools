<template>
  <span>
    <v-container fluid>
      <div v-if="!haveData">
        <v-overlay>
          <v-progress-circular
            indeterminate
            size="64"
          ></v-progress-circular>
        </v-overlay>
      </div>
      <div v-else-if="certificates.length > 0">
        <v-card elevation="2" outlined v-for="(item, i) in certificates" :key="i">
          <v-card-title>
            <v-icon>{{item.icon}}</v-icon>
            <h4 v-text="item.label"/>
          </v-card-title>
          <v-card-text class="text--primary">
            <table>
              <tr><td>Subject</td><td v-text="item.subject"/></tr>
              <tr><td>Valid-From</td><td v-text="item.valid_not_before"/></tr>
              <tr><td>Valid-Until</td><td v-text="item.valid_not_after"/></tr>
            </table>
          </v-card-text>
        </v-card>
      </div>
      <div v-else>
        <p/>
        <v-row align="center" justify="space-around">
          No Certificates available. You may create a new PKI.
        </v-row>
        <p/>
        <div class="text-center">
          <v-btn @click.stop="createPKI()">Create</v-btn>
        </div>
      </div>
    </v-container>
  </span>
</template>

<script>
import axios from "axios";

export default {
  name: "CertificatesView",
  components: {},
  data: function() {
    return { haveData: false, certificates: [] };
  },
  mounted() {
    axios({ method: "GET", url: "/api/certificates" }).then(
      this.handleCertificatesSuccess,this.handleCertificatesFailed);
  },
  methods: {
    createPKI() {
      this.haveData = false;
      axios({ method: "POST", url: "/api/certificates" }).then(
        this.handleCertificatesSuccess, this.handleCertificatesFailed)
    },
    handleCertificatesSuccess(result) {
      var certificates = result.data;
        if (certificates) {
          for (var i = 0; i < certificates.length; i++) {
            certificates[i].uri = "/users/" + certificates[i].label;
            certificates[i].icon = "mdi-certificate";
          }
          this.certificates = certificates;
        } else {
          this.certificates = [];
        }
        this.haveData = true;
    },
    handleCertificatesFailed(error) {
      this.haveData = true;
        if (error != null) {
          error = null;
        }
        this.certificates = []
    }
  }
};
</script>
