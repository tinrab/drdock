<template>
  <main role="main"
        class="root">
    <app-list :items="containers"
              @clickItem="onContainerClick" />
    <app-detail v-if="subject"
                :item="subject" />
  </main>
</template>

<script>
import { mapState } from 'vuex';

import List from '@/components/List';
import Detail from '@/components/Detail';

export default {
  components: {
    'app-list': List,
    'app-detail': Detail,
  },
  computed: mapState({
    containers: (state) => state.containers,
    subject: (state) => state.subject,
  }),
  methods: {
    onContainerClick(id) {
      this.$store.dispatch('selectContainer', id);
    },
  },
};
</script>

<style lang="scss" scoped>
.root {
  display: flex;
}
</style>
