<template>
  <ul class="list-group">
    <li class="list-group-item"
        :class="{active: e.isSelected}"
        v-for="e of items"
        :key="e.id"
        @click="onItemClick(e.id)">
      {{e.name}} {{e.isSelected}}
    </li>
  </ul>
</template>

<script>
export default {
  props: ['items'],
  methods: {
    onItemClick(id) {
      this.$emit('clickItem', id);
    },
  },
};
</script>

<style lang="scss" scoped>
.list-group {
  width: 15rem;
  white-space: nowrap;
  text-overflow: ellipsis;
}
.list-group-item {
  cursor: pointer;
  user-select: none;
  border-radius: 0;
  border-right: 1px solid rgba(0, 0, 0, 0.125);
  border-left: none;
}
</style>
