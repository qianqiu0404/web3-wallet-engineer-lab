<script setup lang="ts">
import { computed, ref } from 'vue'
import { baseline, catalog, failures, type Scenario } from './catalog'

const selectedId = ref(baseline.id)
const revealedSteps = ref(1)
const selected = computed(() => catalog.scenarios.find(item => item.id === selectedId.value) || baseline)

function selectScenario(scenario: Scenario) {
  selectedId.value = scenario.id
  revealedSteps.value = 1
  window.history.replaceState(null, '', `#${scenario.id}`)
}

function revealNext() {
  revealedSteps.value = Math.min(selected.value.steps.length, revealedSteps.value + 1)
}

function reset() {
  revealedSteps.value = 1
}
</script>

<template>
  <main>
    <section class="hero">
      <div class="shell hero-grid">
        <div>
          <p class="eyebrow">Wallet Domain Engine · Technical Evidence</p>
          <h1>Scenario Catalog Inspector</h1>
          <p class="hero-copy">底层 Go 钱包领域引擎的技术检查页，用共享 Catalog 展示资金不变量、故障模型和对应断言；正式交互体验位于独立的 Wallet Reliability Lab。</p>
          <div class="hero-links">
            <a href="https://wallet-reliability-lab.vercel.app" class="primary">打开正式实验台</a>
            <a href="https://github.com/qianqiu0404/web3-wallet-engineer-lab" target="_blank" rel="noopener">查看 Go 测试</a>
          </div>
        </div>
        <aside>
          <span>验证边界</span>
          <strong>确定性模拟 + Go 断言</strong>
          <p>不连接生产服务，不发送链上交易，不保存任何密钥材料。</p>
          <dl><div><dt>场景</dt><dd>1 条基线 / 6 个异常</dd></div><div><dt>事实源</dt><dd>共享 JSON Catalog</dd></div><div><dt>更新</dt><dd>{{ catalog.updatedAt }}</dd></div></dl>
        </aside>
      </div>
    </section>

    <section class="principles">
      <div class="shell principle-grid">
        <div><span>01</span><strong>先确认资金事实</strong><p>链上、订单、冻结、账务分别到了哪一步。</p></div>
        <div><span>02</span><strong>先止损再恢复</strong><p>结果不确定时暂停，不制造第二个资金动作。</p></div>
        <div><span>03</span><strong>恢复必须幂等</strong><p>request_id、交易指纹和 canonical 结果支撑重试。</p></div>
      </div>
    </section>

    <section id="experiments" class="lab shell">
      <header class="section-head">
        <p class="eyebrow">Failure Recovery Experiments</p>
        <h2>检查一条基线与六个异常契约</h2>
        <p>每个实验明确区分故障注入、资金不变量、第一动作、恢复依据与当前边界。</p>
      </header>

      <div class="scenario-tabs" role="tablist" aria-label="实验场景">
        <button :class="{ active: selected.id === baseline.id }" type="button" @click="selectScenario(baseline)"><span>00</span>{{ baseline.title }}</button>
        <button v-for="scenario in failures" :key="scenario.id" :class="{ active: selected.id === scenario.id }" type="button" @click="selectScenario(scenario)"><span>{{ scenario.index }}</span>{{ scenario.title }}</button>
      </div>

      <article class="experiment">
        <header>
          <div><p>{{ selected.kind === 'baseline' ? '正常流程' : `异常实验 ${selected.index}` }}</p><h2>{{ selected.title }}</h2><span>{{ selected.summary }}</span></div>
          <code>{{ selected.goTest }}</code>
        </header>

        <div class="facts-grid">
          <section><small>注入故障</small><p>{{ selected.injectedFault }}</p></section>
          <section><small>资金不变量</small><p>{{ selected.fundInvariant }}</p></section>
          <section><small>第一动作</small><p>{{ selected.firstAction }}</p></section>
          <section><small>恢复依据</small><p>{{ selected.recoveryBasis }}</p></section>
        </div>

        <div class="timeline" aria-live="polite">
          <template v-for="(step, index) in selected.steps" :key="`${selected.id}-${step.state}`">
            <div v-if="index < revealedSteps" class="event" :data-tone="step.tone">
              <span class="event-index">{{ String(index + 1).padStart(2, '0') }}</span>
              <div><small>{{ step.service }} · {{ step.state }}</small><h3>{{ step.title }}</h3><p>{{ step.detail }}</p></div>
            </div>
          </template>
        </div>

        <footer>
          <p><strong>当前边界</strong>{{ selected.currentBoundary }}</p>
          <div><button type="button" @click="reset">重置</button><button v-if="revealedSteps < selected.steps.length" class="primary" type="button" @click="revealNext">显示下一步</button><span v-else>恢复路径已完整显示</span></div>
        </footer>
      </article>
    </section>

    <footer class="site-footer"><div class="shell"><strong>Web3 Wallet Domain Engine</strong><p>这是底层领域模型与契约检查页；交互式讲解由 Wallet Reliability Lab 提供。</p></div></footer>
  </main>
</template>
