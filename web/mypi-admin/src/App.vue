<template>
  <v-app light id="app">
    <app-navigation></app-navigation>
    <v-content transition="slide-x-transition">
        <router-view></router-view>
    </v-content>
    <!--<v-footer app></v-footer>-->
  </v-app>
</template>

<script>
import AppNavigation from '@/components/AppNavigation';

export default {
  name: 'app',
  components: {
    AppNavigation
  },
}
</script>

<style>
</style>
