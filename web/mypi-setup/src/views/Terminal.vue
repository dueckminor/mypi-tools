<template>
  <v-container style="height:90%" fluid>
    <keep-alive>
      <GoTTY
        style="height:100%"
        :path="'/api/hosts/' + $route.params.host + '/terminal/webtty'"
      ></GoTTY>
    </keep-alive>
  </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import GoTTY from "../components/GoTTY";

export default {
  name: "terminal",
  components: {
    GoTTY,
  },
  computed: {},
  created() {
    // eslint-disable-next-line no-console
    console.log("terminal-created");
  },
  mounted() {
    // eslint-disable-next-line no-console
    console.log("terminal-mounted");
  },
  methods: {},
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
