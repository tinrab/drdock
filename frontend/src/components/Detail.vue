<template>
  <div class="root">

    <div class="pb-2">
      <h2>{{item.name}}</h2>
      <table class="table table-sm table-borderless info-table">
        <tbody>
          <tr>
            <th scope="row">ID:</th>
            <td>{{item.id}}</td>
          </tr>
          <tr>
            <th scope="row">Name:</th>
            <td>{{item.name}}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pb-2">
      <h2>Logs</h2>

      <div class="table-responsive">
        <table class="table table-sm">
          <tbody>
            <tr class="log"
                v-for="log of logs"
                :key="log.id">
              <td>{{log.date}}</td>
              <td :style="{textAlign: 'left'}">{{log.message}}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  props: ['item'],
  computed: mapState({
    logs: (state) => state.logs,
  }),
};
</script>

<style lang="scss" scoped>
.root {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
}

.info-table {
  th {
    width: 5rem;
  }
}

.log {
  font-family: SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
}
</style>
