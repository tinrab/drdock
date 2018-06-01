import Vue from 'vue';
import Vuex from 'vuex';

import { Container, Log } from './types';

const MAX_LOGS = 100;
const mutations = {
  ADD_CONTAINER: 'ADD_CONTAINER',
  REMOVE_CONTAINER: 'REMOVE_CONTAINER',
  SELECT_CONTAINER: 'SELECT_CONTAINER',
  ADD_LOGS: 'ADD_LOGS',
};

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    containers: [],
    subject: null,
    logs: [],
  },
  mutations: {
    [mutations.ADD_CONTAINER](state, container) {
      state.containers = [...state.containers, container];
    },
    [mutations.REMOVE_CONTAINER](state, id) {
      state.containers = state.containers.filter((c) => c.id !== id);
    },
    [mutations.SELECT_CONTAINER](state, id) {
      if (state.subject && state.subject.id === id) {
        return;
      }

      for (let c of state.containers) {
        c.isSelected = c.id === id;
        if (c.isSelected) {
          state.subject = c;
        }
      }
      state.logs = [];
    },
    [mutations.ADD_LOGS](state, logs) {
      state.logs = [...logs, ...state.logs].slice(0, MAX_LOGS);
    },
  },
  actions: {
    addContainer({ commit }, container) {
      commit(mutations.ADD_CONTAINER, container);
    },
    removeContainer({ commit }, id) {
      commit(mutations.REMOVE_CONTAINER, id);
    },
    selectContainer({ commit }, id) {
      commit(mutations.SELECT_CONTAINER, id);
    },
    addLogs({ commit }, logs) {
      commit(mutations.ADD_LOGS, logs);
    },
  },
});

store.dispatch('addContainer', new Container('42', 'Test container'));
store.dispatch('addContainer', new Container('56', 'Another one'));
store.dispatch('addContainer', new Container('70', 'Last one'));
for (let i = 1; i <= 3; i++) {
  store.dispatch('addContainer', new Container('100' + i, 'Container #' + i));
}

export default store;
