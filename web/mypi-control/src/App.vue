<template>
  <v-app light id="app">
    <v-app-bar app></v-app-bar>
    <v-main transition="slide-x-transition" width="100%">
      <router-view/>
    </v-main>
    <v-footer app></v-footer>
  </v-app>
</template>

<style></style>
