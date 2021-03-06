
	<template>
    <WorkSpace>
        <template #header>
            <slot name="pageHeader">
                <VHead level="h1">RegionTypeFilter</VHead>
            </slot>
        </template>

        <template #content>
            <slot name="data">
                <table>
                    <thead>
                        <tr>
                            <th v-for="(header, index) in fields" :key="index">{{ header }}</th>
                        </tr>
                    </thead>
            
                    <tbody>
                        <tr
                            v-for="regionTypeFilterItem in regionTypeFilterList"
                            :key="regionTypeFilterItem.Id"
                            @click="selectRegionTypeFilterItem(regionTypeFilterItem)"
                            class="sw-table__row_can-select"
                            :class="{'sw-table__row_is-selected': regionTypeFilterItem.Id === currentRegionTypeFilterItem.item.Id}"
                        >
                            <td v-for="(value, key) in fields" :key="key + '-fields'">
                                <VCheckbox v-if="isCheckbox(regionTypeFilterItem[key])" :checked="regionTypeFilterItem[key]" disabled></VCheckbox>
                                <VText v-else>{{ regionTypeFilterItem[key] }}</VText>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </slot>
            
            <slot name="panel">
                <VPanel
                    v-if="panel.show"
                    width="col3"
                    @close="closePanel"
                >
                    <template #header>
                        <VHead level="h3">{{ panelHeader }}</VHead>
                    </template>
        
                    <template #content>
                        <form @submit.prevent="saveChangesSubmit">
                            <VSet direction="vertical">
                                <VSet
                                    v-for="(filed, key) in editFields" :key="key + '-editFields'"
                                    vertical-align="center"
                                >
                                    <VLabel
                                        width="col4"
                                        :for="`currentRegionTypeFilterItem${key}`"
                                    >{{ filed }}</VLabel>
                                    <VInput
										v-if="isInput(currentRegionTypeFilterItem.item[key])"
                                        v-model="currentRegionTypeFilterItem.item[key]"
                                        width="dyn"
                                        :id="`currentRegionTypeFilterItem${key}`"
                                        @input="changeCurrentRegionTypeFilterItem"
                                    />
									<VCheckbox
										v-if="isCheckbox(currentRegionTypeFilterItem.item[key])"
                                        v-model="currentRegionTypeFilterItem.item[key]"
                                        :id="`currentRegionTypeFilterItem${key}`"
										@input="changeCurrentRegionTypeFilterItem"
									/>
									
                                </VSet>
                            </VSet>
                            <button type="submit" :disabled="!currentRegionTypeFilterItem.hasChange" hidden></button>
                        </form>
                    </template>

                    <template #footer>
                        <VSet>
                            <VButton
                                @click="saveChangesSubmit"
                                accent
                                :text="panelSubmitButtonText"
                                :disabled="!currentRegionTypeFilterItem.hasChange"
                            />
                            <VButton
                                @click="cancelChanges"
                                text="Отменить"
                            />
                        </VSet>
                    </template>
                </VPanel>
            </slot>

            <slot name="confirmationPanel">
                <VPanel
                    v-if="currentRegionTypeFilterItem.showDeleteConfirmation"
                    modal
                    @close="closeConfirmationPanel"
                >
                    <template #content>
                        <VHead level="h3">Удалить элемент?</VHead>
                    </template>

                    <template #footer>
                        <VSet>
                            <VButton
                                text="Подтвердить"
                                accent
                                @click="confirmDeleteHandler"
                            />
                            <VButton
                                text="Отмена"
                                @click="closeConfirmationPanel"
                            />
                        </VSet>
                    </template>
                </VPanel>
            </slot>
        </template>

        <template #footer>
            <slot name="pageFooter">
                <VSet>
                    <VButton
                        v-if="canCreate"
                        text="Добавить"
                        accent
                        @click="showPanel(panel.create)"
                    />
                    <VButton
                        v-if="canDelete"
                        text="Удалить"
                        :disabled="!currentRegionTypeFilterItem.isSelected"
                        @click="deleteRegionTypeFilterItemHandler"
                    />
                </VSet>
            </slot>
        </template>
    </WorkSpace>
</template>

<script>
    import regionTypeFilterData from "../data/RegionTypeFilterData";
    import { RegionTypeFilter } from '../apiModel';
    import { mapGetters, mapMutations, mapActions } from 'vuex';
    import WorkSpace from "swui/src/components/WorkSpace";
    import VHead from "swui/src/components/VHead";
    import VSet from "swui/src/components/VSet";
    import VLabel from "swui/src/components/VLabel";
    import VInput from "swui/src/components/VInput";
    import VCheckbox from "swui/src/components/VCheckbox";
    import VText from "swui/src/components/VText";
    import VPanel from "swui/src/components/VPanel";
    import VButton from "swui/src/components/VButton";

    export default {
        name: 'RegionTypeFilterGen',

        components: {VButton, VPanel, VText, VInput, VLabel, VSet, VHead, WorkSpace, VCheckbox},

        props: {
            fields: {
                type: Object,
                default() {
                    const regionTypeFilterItem = new RegionTypeFilter();
                    const fieldsObj = {};

                    for (let prop in regionTypeFilterItem) {

                        if (regionTypeFilterItem.hasOwnProperty(prop)) {
                            fieldsObj[prop] = prop;
                        }

                    }

                    return fieldsObj;
                }
            },
            editFields: {
                type: Object,
                default() {
                    const regionTypeFilterItem = new RegionTypeFilter();
                    const fieldsObj = {};

                    for (let prop in regionTypeFilterItem) {

                        if (regionTypeFilterItem.hasOwnProperty(prop)) {
                            fieldsObj[prop] = prop;
                        }

                    }

                    return fieldsObj;
                }
            },
            canDelete: {
                type: Boolean,
                default: true,
            },
            canCreate: {
                type: Boolean,
                default: true,
            },
        },

        data() {
            return regionTypeFilterData;
        },

        created() {
			this.onCreated();
        },

        computed: {
            ...mapGetters({
                regionTypeFilterList: 'getListRegionTypeFilter'
            }),
            isPanelCreate() {
                return this.panel.type === this.panel.create;
            },
            isPanelEdit() {
                return this.panel.type === this.panel.edit;
            },
            panelHeader() {
                if (this.isPanelCreate) {
                    return this.panelHeaderCreate;
                }

                if (this.isPanelEdit) {
                    return this.panelHeaderEdit;
                }

                return  '';
            },
            panelSubmitButtonText() {
                if (this.isPanelCreate) {
                    return this.panelSubmitButtonTextCreate;
                }

                if (this.isPanelEdit) {
                    return this.panelSubmitButtonTextEdit;
                }

                return  '';
            },
            isCheckbox() {
                return data => {
                    return typeof data === "boolean";
                }
            },
            isInput() {
                return data => {
                    return ! this.isCheckbox(data);
                }
            },
        },

        methods: {
            ...mapActions([
                'findRegionTypeFilter',
                'updateRegionTypeFilter',
                'deleteRegionTypeFilter',
                'createRegionTypeFilter',
            ]),

            ...mapMutations([
                'addRegionTypeFilterItemToList',
                'deleteRegionTypeFilterFromList',
                'updateRegionTypeFilterById',
            ]),

			onCreated() {
				this.fillRegionTypeFilterFilter();
	            this.fetchRegionTypeFilterData();
			},

            fillRegionTypeFilterFilter() {
                this.regionTypeFilterFilter.CurrentPage = 1;
                this.regionTypeFilterFilter.PerPage = 1000;
            },

            fetchRegionTypeFilterData() {
                return this.findRegionTypeFilter({
                    filter: this.regionTypeFilterFilter
                });
            },

            /**
             *
             * @param type
             */
            showPanel(type) {
                if (type === this.panel.create) {
                    this.panel.type = this.panel.create;
                    this.clearPanelRegionTypeFilterItem();
                } else if (type === this.panel.edit) {
                    this.panel.type = this.panel.edit;
                    this.currentRegionTypeFilterItem.isSelected = true;
                }

                this.panel.show = true;
            },

            closePanel() {
                this.panel.show = false;
                this.currentRegionTypeFilterItem.isSelected = false;
                this.clearPanelRegionTypeFilterItem();
            },

            selectRegionTypeFilterItem(regionTypeFilterItem) {
                this.showPanel(this.panel.edit);
                this.currentRegionTypeFilterItem.isSelected = true;
                Object.assign(this.currentRegionTypeFilterItem.item, regionTypeFilterItem);
            },

            changeCurrentRegionTypeFilterItem() {
                this.currentRegionTypeFilterItem.hasChange = true;
            },

            cancelChanges() {
                this.clearPanelRegionTypeFilterItem();
                this.closePanel();
            },

            clearPanelRegionTypeFilterItem() {
                this.currentRegionTypeFilterItem.item = new RegionTypeFilter();
                this.currentRegionTypeFilterItem.hasChange = false;
            },

            saveChangesSubmit() {
                if (this.isPanelCreate) {
                    this.createRegionTypeFilterItemSubmit();
                    return;
                }

                if (this.isPanelEdit) {
                    this.editRegionTypeFilterItemSubmit();
                }
            },

            createRegionTypeFilterItemSubmit() {
                return this.createRegionTypeFilter({
					data: this.currentRegionTypeFilterItem.item,
                }).then((response) => {

                    if (response.Model) {
                        this.addRegionTypeFilterItemToList(response.Model);
                        this.clearPanelRegionTypeFilterItem();
                    } else {
                        console.error('Ошибка создания записи: ', response.Error);
                    }

                }).catch(error => {
                    console.error('Ошибка создания записи: ', error);
                });
            },

            editRegionTypeFilterItemSubmit() {

                if (this.currentRegionTypeFilterItem.hasChange) {
                    return this.updateRegionTypeFilter({
                        id: this.currentRegionTypeFilterItem.item.Id,
                        data: this.currentRegionTypeFilterItem.item,
                    }).then((response) => {

                        if (response.Model) {
                            this.updateRegionTypeFilterById(response.Model);
                            this.currentRegionTypeFilterItem.hasChange = false;
                            this.clearPanelRegionTypeFilterItem();
                            this.closePanel();
                        } else {
                            console.error('Ошибка изменения записи: ', response.Error);
                        }

                    }).catch(error => {
                        console.error('Ошибка изменения записи: ', error);
                    });

                } else {
					return new Promise(function(resolve, reject) {reject("Item has no changes. Nothing to save");})
				}
            },

            deleteRegionTypeFilterItemHandler() {
                let deletedItemId = this.currentRegionTypeFilterItem.item.Id;

                if (!this.currentRegionTypeFilterItem.canDelete) {
                    this.currentRegionTypeFilterItem.showDeleteConfirmation = true;
                    return;
                }

                this.deleteRegionTypeFilter({
                    id: deletedItemId
                }).then(response => {

                    if (response.IsSuccess) {
                        this.deleteRegionTypeFilterFromList(deletedItemId);
                        this.clearPanelRegionTypeFilterItem();
                        this.currentRegionTypeFilterItem.canDelete = false;
                        this.currentRegionTypeFilterItem.isSelected = false;
                        this.panel.show = false;
                    } else {
                        console.error('Ошибка удаления элемента: ', response.Error);
                    }

                }).catch(error => {
                    console.error('Ошибка удаления элемента: ', error);
                });
            },

            confirmDeleteHandler() {
                this.currentRegionTypeFilterItem.showDeleteConfirmation = false;
                this.currentRegionTypeFilterItem.canDelete = true;
                this.deleteRegionTypeFilterItemHandler();
            },

            closeConfirmationPanel() {
                this.currentRegionTypeFilterItem.showDeleteConfirmation = false;
            },
        },
    }
</script>

<style lang="scss">

</style>
