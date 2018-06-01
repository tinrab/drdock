<template>
  <div class="root">

    <nav class="navbar navbar-dark bg-dark">
      <a href="/"
         class="navbar-brand">Dr. Dock</a>
    </nav>

    <app-dashboard />

  </div>
</template>

<script>
import Dashboard from '@/components/Dashboard';

export default {
  components: {
    'app-dashboard': Dashboard,
  },
};
</script>

<style lang="scss">
html,
body {
  height: 100%;
}
</style>

<style lang="scss" scoped>
.root {
  height: 100%;
}
</style>
