<template>
  <v-app light id="app">
    <v-navigation-drawer app v-model="drawer" hide-overlay stateless>
      <v-list nav>
        <v-list-item v-for="item in items" :key="item" :title="item.title" :prepend-icon="item.icon" :to="item.url"/>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>MYPI-Setup</v-toolbar-title>
      <v-breadcrumbs divider=">"> </v-breadcrumbs>
    </v-app-bar>
    <v-main transition="slide-x-transition" width="100%">
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: "app",
  components: {
  },
  created() {
    // eslint-disable-next-line no-console
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
          title: "Certificates",
          icon: "mdi-certificate-outline",
          url: "/certificates",
        },
        { 
          title: "Hosts", 
          icon: "mdi-server", 
          url: "/hosts" },
      ],
    };
  },
  methods: {},
};
</script>

<style></style>
