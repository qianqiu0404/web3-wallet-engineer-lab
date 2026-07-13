import { describe, expect, it } from 'vitest'
import { baseline, catalog, failures } from './catalog'

describe('shared reliability catalog', () => {
  it('contains one baseline and six failure experiments', () => {
    expect(baseline.index).toBe('00')
    expect(failures).toHaveLength(6)
    expect(catalog.scenarios).toHaveLength(7)
  })

  it('keeps recovery evidence complete and relations unique', () => {
    const ids = new Set<string>()
    const relations = new Set<string>()
    for (const scenario of catalog.scenarios) {
      expect(ids.has(scenario.id)).toBe(false)
      ids.add(scenario.id)
      expect(scenario.fundInvariant).toBeTruthy()
      expect(scenario.firstAction).toBeTruthy()
      expect(scenario.recoveryBasis).toBeTruthy()
      expect(scenario.goTest).toMatch(/^Test/)
      expect(scenario.steps.length).toBeGreaterThanOrEqual(4)
      if (scenario.kind === 'failure') {
        expect(scenario.failureSlug).toBeTruthy()
        expect(relations.has(scenario.failureSlug)).toBe(false)
        relations.add(scenario.failureSlug)
      }
    }
  })
})
