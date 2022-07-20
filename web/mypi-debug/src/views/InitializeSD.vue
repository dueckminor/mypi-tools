<template>
  <v-container style="height:90%" fluid>
    <GoTTY style="height:100%" path="/api/hosts/localhost/actions/makesd/webtty"></GoTTY>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import GoTTY from "../components/GoTTY";
import axios from "axios";

export default {
  name: "initialize_sd",
  components: {
    GoTTY
  },
  computed: {},
  mounted() {
    axios.post("/api/hosts/localhost/actions/makesd", {
      Disk:          this.$route.params.Disk,
      Hostname:      this.$route.params.Hostname,
      AlpineVersion: this.$route.params.AlpineVersion,
      AlpineArch:    this.$route.params.AlpineArch,
    })
  },
  methods: {}
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
