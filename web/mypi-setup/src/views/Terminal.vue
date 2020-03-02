<template>
  <v-container style="height:90%" fluid>
    <GoTTY style="height:100%" :path="'/api/mypi/'+$route.params.host+'/terminal'"></GoTTY>
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
