<template>
  <v-container fluid> </v-container>
</template>
<!-- ----------------------------------------------------------------------- -->
<script>
import axios from "axios";

export default {
  name: "initialize_sd",
  components: {},
  computed: {},
  mounted() {
    axios.get("/api/hosts");
  },
  methods: {},
};
</script>
<!-- ----------------------------------------------------------------------- -->
<style scoped></style>
