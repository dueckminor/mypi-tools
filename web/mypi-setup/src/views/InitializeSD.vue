<template>
  <v-container style="height:90%" fluid>
    <GoTTY style="height:100%" path="/api/hosts/localhost/actions/makesd/webtty"></GoTTY>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import GoTTY from "../components/GoTTY";

export default {
  name: "initialize_sd",
  components: {
    GoTTY
  },
  computed: {},
  mounted() {},
  methods: {}
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
