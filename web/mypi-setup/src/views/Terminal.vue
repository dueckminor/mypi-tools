<template>
  <v-container style="height:90%" fluid>
    <GoTTY style="height:100%" :path="'/api/hosts/'+$route.params.host+'/terminal/webtty'"></GoTTY>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import GoTTY from "../components/GoTTY";

export default {
  name: "terminal",
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
