<template>
  <v-app light id="app">
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

<script>
import { defineComponent } from 'vue'
import { useServiceStore } from "@/stores/ServiceStore";


export default defineComponent({
  name: "app",
  sockets: {
    connect() {
      console.log('socket connected')
    },
    customEmit() {
      console.log('this method was fired by the socket server. eg: io.emit("customEmit", data)')
    }
  },
  setup() {
    const serviceStore = useServiceStore()
    console.log("serviceStore.fill")
    serviceStore.fill()
  },
  components: {
  },
  created() {
    // eslint-disable-next-line no-console
    //console.log(serviceStore);
    console.log("app-created");
  },
  mounted() {
    // eslint-disable-next-line no-console
    console.log("app-mounted");
  },
  data: function() {
    return {
      showLogin: true,
      drawer: null,
      items: [
        {
          title: "Debug",
          url: "/debug",
        },
      ],
    };
  },
  methods: {},
});
</script>

<style></style>
