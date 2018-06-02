import Vue from 'vue';
import Vuex from 'vuex';

import { Container, Log } from './types';
import Socket from './socket';

const MAX_LOGS = 100;
const WS_URL = 'ws://localhost:8000/ws';

const mutations = {
  ADD_CONTAINER: 'ADD_CONTAINER',
  REMOVE_CONTAINER: 'REMOVE_CONTAINER',
  SELECT_CONTAINER: 'SELECT_CONTAINER',
  FETCH_CONTAINERS: 'FETCH_CONTAINERS',
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

      const containers = state.containers;
      for (let c of containers) {
        c.isSelected = c.id === id;
        if (c.isSelected) {
          state.subject = c;
        }
      }
      state.logs = [];
    },
    [mutations.FETCH_CONTAINERS](state, containers) {
      state.containers = containers;
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
    fetchContainers({ commit }, containers) {
      commit(mutations.FETCH_CONTAINERS, containers);
    },
  },
});

const socket = new Socket(WS_URL);
socket.onAddContainer = (c) => store.dispatch('addContainer', c);
socket.onRemoveContainer = (id) => store.dispatch('removeContainer', id);
socket.onFetchContainers = (containers) => store.dispatch('fetchContainers', containers);
socket.open();

export default store;
