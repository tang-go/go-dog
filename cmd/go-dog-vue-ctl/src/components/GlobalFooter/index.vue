<template>
  <global-footer class="footer custom-render">
    <template v-slot:links>
      <a href="" target="_blank">go-dog</a>
      <a href="https://github.com/tang-go/go-dog" target="_blank">Github</a>
      <a href="" target="_blank">@TangJie</a>
    </template>
    <template v-slot:copyright>
      <a href="" target="_blank">服务治理中心</a>
    </template>
  </global-footer>
</template>

<script>
import { GlobalFooter } from '@ant-design-vue/pro-layout'

export default {
  name: 'ProGlobalFooter',
  components: {
    GlobalFooter
  }
}
</script>
