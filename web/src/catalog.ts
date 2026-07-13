import source from '../../scenarios/catalog.json'

export type ScenarioKind = 'baseline' | 'failure'
export type StepTone = 'active' | 'success' | 'warning' | 'danger' | 'recovered'

export interface ScenarioStep {
  service: string
  state: string
  title: string
  detail: string
  tone: StepTone
}

export interface Scenario {
  id: string
  index: string
  title: string
  kind: ScenarioKind
  failureSlug: string
  summary: string
  injectedFault: string
  fundInvariant: string
  firstAction: string
  recoveryBasis: string
  currentBoundary: string
  goTest: string
  steps: ScenarioStep[]
}

export interface ScenarioCatalog {
  version: number
  updatedAt: string
  scenarios: Scenario[]
}

export const catalog = source as ScenarioCatalog
export const baseline = catalog.scenarios.find(item => item.kind === 'baseline')!
export const failures = catalog.scenarios.filter(item => item.kind === 'failure')
