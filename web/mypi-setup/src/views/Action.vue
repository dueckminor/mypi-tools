<template>
  <v-container style="height:90%" fluid>
    <GoTTY
      style="height:100%"
      :path="'/api/hosts/' + $route.params.host + '/actions/' + $route.params.action + '/webtty'"
    ></GoTTY>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import GoTTY from "../components/GoTTY";
//import axios from "axios";

export default {
  name: "action",
  components: {
    GoTTY
  },
  computed: {},
  mounted() {
  },
  methods: {}
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
