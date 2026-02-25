const states = ['DRAFT', 'SUBMITTED', 'REJECTED', 'REVIEWED', 'APPROVED', 'REGISTERED', 'ARCHIVED'] as const
export type ContractTemplateState = typeof states[number]
export const contractTemplateStates: ContractTemplateState[] = [...states]
