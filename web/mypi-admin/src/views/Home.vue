<template>
    <span>
    <v-list shaped>
      <v-subheader>REPORTS</v-subheader>
      <v-list-item-group v-model="item" color="primary">
        <v-list-item
          v-for="(item, i) in items"
          :key="i"
        >
          <v-list-item-icon>
            <v-icon v-text="item.icon"></v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title v-text="item.text"></v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
    </v-list>
    </span>
</template>

<script>
// import HomeHero from '@/components/HomeHero';
// import HomeDetails from '@/components/HomeDetails';
// import HomePlans from '@/components/HomePlans';

export default {
    name: 'home',
    components: {
        // HomeHero,
        // HomeDetails,
        // HomePlans
    }
};
</script>
