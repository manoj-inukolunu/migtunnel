const tunnelList = {
  props: [], data() {
    return {list: undefined}
  }, mounted() {
    fetch('/list').then(res => res.json()).then(d => {
      this.list = d;
    })
  }, template: `
      <a href="#" v-for="tunnel in list" v-bind:key="tunnel.HostName" class="block py-2.5 px-4 rounded transition duration-200 hover:bg-blue-700 hover:text-white">
        <span>{{ tunnel.HostName }}</span>
      </a>
  `
}