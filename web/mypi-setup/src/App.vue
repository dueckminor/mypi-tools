<template>
  <v-app light id="app">
    <v-app-bar app></v-app-bar>
    <v-content transition="slide-x-transition" width="100%">
        <router-view></router-view>
    </v-content>
    <v-footer app></v-footer>
  </v-app>
</template>

<script>

export default {
  name: "app",
  components: {
  },
  data: function() {
    return {
      showLogin: true
    };
  },
  methods: {
  }
};
</script>

<style></style>
