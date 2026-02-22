describe('Navigation', () => {
  beforeEach(() => {
    cy.visit('/')
  })

  it('shows home page with title and nav links', () => {
    cy.contains('h1', 'Music Service')
    cy.contains('a', 'Albums')
    cy.contains('a', 'Add Album')
  })

  it('navigates to Albums page', () => {
    cy.contains('a', 'Albums').click()
    cy.url().should('include', '#/albums')
    cy.contains('h2', 'Albums')
  })

  it('navigates to Add Album page', () => {
    cy.contains('a', 'Add Album').click()
    cy.url().should('include', '#/add')
    cy.contains('h2', 'Add Album')
  })

  it('navigates back to home via Music Service link', () => {
    cy.contains('a', 'Add Album').click()
    cy.contains('a', 'Music Service').click()
    cy.url().should('match', /#\/?$/)
  })
})
