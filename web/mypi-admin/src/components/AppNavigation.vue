<template>
  <span>
    <v-app-bar app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <!--<v-app-bar-nav-icon @click="updateAppDrawer"> </v-app-bar-nav-icon>-->
      <v-toolbar-title>Foo</v-toolbar-title>
        <v-breadcrumbs divider=">">
  </v-breadcrumbs>
    </v-app-bar>
    <v-navigation-drawer
      app
      v-model="drawer"
      hide-overlay
      stateless
    >
      <v-list nav>
        <template v-for="(item, index) in items">
          <v-list-item link :key="index" :to="item.url">
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>

            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>
  </span></template
>

<script>
export default {
  name: "AppNavigation",
  data() {
    return {
      appTitle: "Meal Prep",
      drawer: null,
      items: [
        { title: "Domains", icon: "mdi-domain", url: "/domains" },
        { title: "Users", icon: "mdi-account-multiple", url: "/users" },
        { title: "Apps", icon: "mdi-apps", url: "/apps" },
        { title: "Docker", icon: "mdi-docker", url: "/docker" },
        { title: "Settings", icon: "mdi-settings", url: "/settings" }
      ]
    };
  },
  computed: {
    isAuthenticated() {
      return true;
      // return this.$store.getters.isAuthenticated;
    }
  },
  methods: {
    updateAppDrawer() {
      this.drawer = true;
    },
    logout() {
      //this.$store.dispatch('userSignOut');
    }
  }
};
</script>

<style scoped>
a {
  color: white;
  text-decoration: none;
}
</style>
