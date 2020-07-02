<template>
    <div class="diractory-layout">
        <bk-input class="dir-search"
            v-model.trim="dirSearch"
            clearable
            :placeholder="$t('分组目录')">
        </bk-input>
        <div class="dir-header"
            :class="{ 'active': acitveDirId === null }"
            @click="handleResourceClick">
            <span class="title">{{$t('资源池')}}</span>
            <i class="icon-cc-plus" v-bk-tooltips.top="$t('新建目录')" @click.stop="handleShowCreate"></i>
        </div>
        <ul class="dir-list">
            <li class="dir-item edit-status" v-if="createInfo.active">
                <bk-input
                    ref="createdDir"
                    v-click-outside="handleCancelCreate"
                    style="width: 100%"
                    v-validate="'required|singlechar|length:256'"
                    :placeholder="$t('请输入目录名称，回车结束')"
                    v-model="createInfo.name"
                    @enter="handleConfirm(true)">
                </bk-input>
            </li>
            <li v-for="(dir, index) in filterDirList"
                :key="index"
                :class="{
                    'dir-item': true,
                    'edit-status': editDir.id === dir.bk_module_id,
                    'selected': acitveDirId === dir.bk_module_id
                }"
                @click="handleSearchHost(dir)">
                <template v-if="editDir.id === dir.bk_module_id">
                    <bk-input
                        style="width: 100%"
                        v-click-outside="handleCancelEdit"
                        v-validate="'required|singlechar|length:256'"
                        :placeholder="$t('请输入目录名称，回车结束')"
                        :ref="`dir-node-${dir.bk_module_id}`"
                        v-model="editDir.name"
                        @enter="handleConfirm(false)"
                        @click.native.stop>
                    </bk-input>
                </template>
                <template v-else>
                    <i class="icon-cc-memory"></i>
                    <span class="dir-name" :title="dir.bk_module_name">{{dir.bk_module_name}}</span>
                    <cmdb-dot-menu class="dir-operation" color="#3A84FF" @click.native.stop="handleCloseInput">
                        <div class="dot-content">
                            <bk-button
                                class="menu-btn"
                                :text="true"
                                @click="handleResetName(dir)">
                                {{$t('重命名')}}
                            </bk-button>
                            <bk-button
                                class="menu-btn"
                                :text="true"
                                v-bk-tooltips.right="{
                                    content: $t('主机不为空，不能删除'),
                                    disabled: !dir.host_count
                                }"
                                @click="handleDelete(dir, index)">
                                {{$t('删除')}}
                            </bk-button>
                        </div>
                    </cmdb-dot-menu>
                    <span class="host-count">{{dir.host_count}}</span>
                </template>
            </li>
        </ul>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import Bus from '@/utils/bus.js'
    import RouterQuery from '@/router/query'
    export default {
        data () {
            return {
                dirSearch: '',
                resetName: false,
                createInfo: {
                    active: false,
                    name: ''
                },
                editDir: {
                    id: null,
                    name: ''
                },
                acitveDirId: null
            }
        },
        computed: {
            ...mapGetters('resourceHost', [
                'directoryList'
            ]),
            filterDirList () {
                if (this.dirSearch) {
                    return this.directoryList.filter(module => module.bk_module_name.indexOf(this.dirSearch) > -1)
                }
                return this.directoryList
            }
        },
        watch: {
            acitveDirId (id) {
                RouterQuery.set({
                    directory: id,
                    page: 1,
                    _t: Date.now()
                })
            }
        },
        async created () {
            Bus.$on('refresh-dir-count', this.getDirectoryList)
            this.getDirectoryList()
        },
        beforeDestroy () {
            Bus.$off('refresh-dir-count', this.getDirectoryList)
        },
        methods: {
            async getDirectoryList () {
                try {
                    const { info } = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
                        params: {
                            page: {
                                sort: '-bk_module_name'
                            }
                        },
                        config: {
                            requestId: 'getDirectoryList'
                        }
                    })
                    this.$store.commit('resourceHost/setDirectoryList', info)
                } catch (error) {
                    console.error(error)
                }
            },
            async createdDir () {
                try {
                    const data = await this.$store.dispatch('resourceDirectory/createDirectory', {
                        params: {
                            bk_module_name: this.createInfo.name
                        }
                    })
                    const newDir = {
                        bk_module_id: data.created.id,
                        bk_module_name: this.createInfo.name,
                        host_count: 0
                    }
                    this.$store.commit('resourceHost/addDirectory', newDir)
                    this.$success(this.$t('新建成功'))
                    this.handleCancelCreate()
                } catch (e) {
                    console.error(e)
                }
            },
            async updateDir () {
                try {
                    await this.$store.dispatch('resourceDirectory/updateDirectory', {
                        moduleId: this.editDir.id,
                        params: {
                            bk_module_name: this.editDir.name
                        },
                        config: {
                            requestId: 'updateDir'
                        }
                    })
                    const target = this.directoryList.find(dir => dir.bk_module_id === this.editDir.id)
                    this.$store.commit('resourceHost/updateDirectory', Object.assign({}, target, {
                        bk_module_id: this.editDir.id,
                        bk_module_name: this.editDir.name
                    }))
                    this.$success(this.$t('修改成功'))
                    this.handleCancelEdit()
                } catch (e) {
                    console.error(e)
                }
            },
            handleSearchHost (active = {}) {
                this.$store.commit('resourceHost/setActiveDirectory', active)
                this.acitveDirId = active.bk_module_id
                Bus.$emit('refresh-resource-list')
            },
            handleResourceClick () {
                this.$store.commit('resourceHost/setActiveDirectory', null)
                this.acitveDirId = null
                Bus.$emit('refresh-resource-list')
            },
            handleCancelCreate () {
                this.createInfo.active = false
                this.createInfo.name = ''
            },
            handleShowCreate () {
                this.createInfo.active = true
                this.$nextTick(() => {
                    this.$refs.createdDir.$refs.input.focus()
                })
            },
            handleCancelEdit () {
                this.editDir.id = null
                this.editDir.name = ''
            },
            async handleConfirm (isCreate) {
                if (!await this.$validator.validateAll()) {
                    this.$error(this.$t('请正确目录名称'))
                    return
                }
                if (isCreate) {
                    this.createdDir()
                } else {
                    this.updateDir()
                }
            },
            handleResetName (dir) {
                this.editDir.id = dir.bk_module_id
                this.editDir.name = dir.bk_module_name
                this.$nextTick(() => {
                    this.$refs[`dir-node-${dir.bk_module_id}`][0].$refs.input.focus()
                })
            },
            handleCloseInput () {
                this.handleCancelCreate()
                this.handleCancelEdit()
            },
            async handleDelete (dir, index) {
                if (dir.host_count) {
                    this.$error(this.$t('目标包含主机, 不允许删除'))
                    return
                }
                this.$bkInfo({
                    title: this.$t('确认确定删除目录'),
                    subTitle: this.$t('即将删除目录', { name: dir.bk_module_name }),
                    extCls: 'bk-dialog-sub-header-center',
                    confirmFn: async () => {
                        try {
                            await this.$store.dispatch('resourceDirectory/deleteDirectory', {
                                moduleId: dir.bk_module_id,
                                config: {
                                    requestId: 'deleteDirectory'
                                }
                            })
                            if (dir.bk_module_id === this.acitveDirId) {
                                this.acitveDirId = this.directoryList[index - 1].bk_module_id
                            }
                            this.$store.commit('resourceHost/deleteDirectory', dir.bk_module_id)
                            this.$success(this.$t('删除成功'))
                        } catch (e) {
                            console.error(e)
                        }
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .diractory-layout {
        height: 100%;
        overflow: hidden;
        .dir-search {
            display: block;
            width: auto;
            margin: 18px 20px 14px;
        }
        .dir-header {
            @include space-between;
            padding: 0 20px;
            height: 42px;
            line-height: 42px;
            background-color: #F0F1F5;
            cursor: pointer;
            &:hover,
            &.active {
                background-color: #E1ECFF;
                .icon-cc-plus {
                    background-color: #3A84FF;
                }
            }
            .title {
                font-weight: bold;
                font-size: 14px;
            }
            .icon-cc-plus {
                width: 18px;
                height: 18px;
                line-height: 18px;
                text-align: center;
                color: #FFFFFF;
                background-color: #C4C6CC;
                border-radius: 2px;
                cursor: pointer;
            }
        }
        .dir-list {
            height: calc(100% - 106px);
            padding-bottom: 10px;
            @include scrollbar-y;
        }
        .dir-item {
            display: flex;
            align-items: center;
            height: 36px;
            padding: 0 20px;
            margin: 6px 0;
            cursor: pointer;
            &:first-child {
                margin-top: 0;
            }
            &:not(.edit-status):not(.disabled):hover,
            &:not(.edit-status).selected {
                background-color: #E1ECFF;
                .icon-cc-memory {
                    color: #3A84FF;
                }
                .dir-name {
                    color: #3A84FF;
                }
                .dir-operation {
                    display: block;
                    opacity: 1;
                }
                .host-count {
                    color: #FFFFFF;
                    background-color: #A2C5FD;
                }
            }
            &.disabled {
                .icon-cc-memory {
                    color: #DCDEE5 !important;
                }
                .dir-name {
                    color: #C4C6CC;
                }
            }
            .icon-cc-memory {
                font-size: 16px;
                margin-right: 10px;
                color: #C4C6CC;
            }
            .dir-name {
                flex: 1;
                font-size: 14px;
                color: #63656E;
                @include ellipsis;
            }
            .dir-operation {
                width: 20px;
                margin-right: 8px;
                opacity: 0;
            }
            .host-count {
                height: 18px;
                line-height: 17px;
                font-size: 12px;
                padding: 0 5px;
                color: #979BA5;
                text-align: center;
                background-color: #F0F1F5;
                border-radius: 2px;
            }
            .icon-cc-lock {
                font-size: 14px;
                color: #C4C6CC;
            }
        }
    }
    .dot-content {
        width: 90px;
        padding: 6px 0;
        .auth-box {
            display: block;
        }
        .menu-btn {
            display: block;
            width: 100%;
            height: 32px;
            line-height: 32px;
            padding: 0 8px;
            text-align: left;
            color: #63656E;
            outline: none;
            &:hover {
                color: #3A84FF;
                background-color: #E1ECFF;
            }
            &:disabled {
                color: #DCDEE5;
                background-color: transparent;
            }
            &.no-allow-btn {
                cursor: not-allowed;
                color: #DCDEE5;
                background-color: transparent;
            }
        }
    }
</style>