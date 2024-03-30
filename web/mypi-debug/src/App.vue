<template>
  <v-app id="app">
    <v-navigation-drawer app v-model="drawer" hide-overlay stateless>
      <v-list nav>
        <v-list-item v-for="item in items" :key="item" :title="item.title" :prepend-icon="item.icon" :to="item.url"/>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>mypi-debug</v-toolbar-title>
      <v-breadcrumbs divider=">"> </v-breadcrumbs>
    </v-app-bar>
    <v-main transition="slide-x-transition" width="100%">
      <router-view/>
    </v-main>
  </v-app>
</template>

<script setup>
import { useServiceStore } from "@/stores/ServiceStore";

const serviceStore = useServiceStore()
console.log("serviceStore.fill")
serviceStore.fill()
</script>
