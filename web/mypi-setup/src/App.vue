<template>
  <v-app light id="app">
    <app-navigation></app-navigation>
    <!--
    <v-app-bar app
      ><v-app-bar-nav-icon></v-app-bar-nav-icon
      ><v-toolbar-title>MYPI-Setup</v-toolbar-title></v-app-bar
    >
    -->
    <v-main transition="slide-x-transition" width="100%">
      <router-view></router-view>
    </v-main>
  </v-app>
</template>

<script>
import AppNavigation from "@/components/AppNavigation";

export default {
  name: "app",
  components: {
    AppNavigation,
  },
  data: function() {
    return {
      showLogin: true,
    };
  },
  methods: {},
};
</script>

<style></style>
