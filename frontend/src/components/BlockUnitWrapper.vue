<template>
  <div>
    <q-expansion-item
      v-model="expandByDefault"
      :dense="true"
      :expand-separator="true"
    >
      <template #header>
        <div class="full-width row items-center justify-between">
          <div>
            <slot name="header-extend-left"></slot>
          </div>

          <div class="row items-center">
            <slot name="header-extend-right"></slot>
            <q-separator spaced style="width: 2px" vertical />
            <template v-if="opIndex === 0">
              <teleport
                v-if="enableTeleport"
                :to="'#toggle-condition-block_' + scope.blockIndex"
              >
                <q-toggle
                  size="sm"
                  v-model="enableCondition"
                  left-label
                  label="Condition"
                  class="text-bold"
                  color="green"
                />
              </teleport>
              <q-toggle
                v-else
                size="sm"
                v-model="enableCondition"
                left-label
                label="Condition"
                class="text-bold"
                color="green"
              />
            </template>
          </div>
        </div>

        <slot name="head"> </slot>
      </template>
      <div>
        <slot></slot>
      </div>
    </q-expansion-item>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { BlockScope } from "..";
import { defaultCondition } from "~/utils/utills";
import { useInject } from "~/composable/useProviderInjector";
const props = defineProps<{
  scope: BlockScope;
  hideToggleBtn: boolean;
  expandable?: boolean;
}>();
const enableCondition = ref(false);
const expandByDefault = ref(props.expandable ?? false);
const enableTeleport = ref(false);

const operationIndex = useInject("operationIndex");
const opIndex = ref();
onMounted(() => {
  opIndex.value = operationIndex?.();
});

watch(
  () => enableCondition.value,
  () => {
    const op = props.scope.operation;
    if (!op || op?.condition == undefined) return;
    op.condition.disableCondition = !enableCondition.value; // keep the data, just switch it off
    nextTick().then(() => {
      setTimeout(() => {
        enableTeleport.value = !enableTeleport.value;
      }, 500);
    });
    if (op.data.length < 2) {
      return;
    }
    op.data.forEach((d, i) => {
      props.scope.operations?.push({
        condition: { ...defaultCondition },
        data: [d],
      });
      op.data.splice(i, 1);
    });
  },
);
</script>

<style></style>
