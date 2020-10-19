<template>
  <div>
    <div class="table-operator">
      <a-button type="primary" icon="plus" @click="addRole">新建</a-button>
    </div>
    <s-table ref="table" size="default" :columns="columns" :data="loadData">
      <span slot="action" >
        <a>修改</a>
        <a-divider type="vertical"/>
        <a>删除</a>
      </span>
    </s-table>

    <a-modal v-model="addModel" title="添加角色" @ok="addOk">
      <p>Some contents...</p>
      <p>Some contents...</p>
      <p>Some contents...</p>
    </a-modal>
  </div>
</template>

<script>
  import { STable } from '@/components'
  import { getRoleList } from '@/api/manage'

  export default {
    components: {
      STable
    },
    data () {
      return {
        addModel: false,
        columns: [
          {
            title: '角色ID',
            dataIndex: 'roleId'
          },
          {
            title: '角色名称',
            dataIndex: 'name'
          },
          {
            title: '描述',
            dataIndex: 'description'
          },
          {
            title: '创建时间',
            dataIndex: 'time'
          },
          {
            table: '操作',
            dataIndex: 'action',
            scopedSlots: { customRender: 'action' }
          }
        ],
        // 查询条件参数
        queryParam: {},
        // 加载数据方法 必须为 Promise 对象
        loadData: parameter => {
          console.log('test', Object.assign(parameter, this.queryParam))
          return getRoleList(Object.assign(parameter, this.queryParam)).then(res => {
            return res.body
          })
        }
      }
    },
    methods: {
      addRole () {
        this.addModel = true
      },
      addOk () {
        this.addModel = false
      },
      addCancel () {
        this.addModel = false
      },
      edit (row) {
        // axios 发送数据到后端 修改数据成功后
        // 调用 refresh() 重新加载列表数据
        // 这里 setTimeout 模拟发起请求的网络延迟..
        setTimeout(() => {
          this.$refs.table.refresh() // refresh() 不传参默认值 false 不刷新到分页第一页
        }, 1500)
      }
    }
  }
</script>
