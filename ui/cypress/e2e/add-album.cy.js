describe('Add Album', () => {
  beforeEach(() => {
    cy.intercept('POST', '**/api/v1/album', {
      statusCode: 201,
      body: { id: 1, title: 'Blue Train', artist: 'John Coltrane', price: 56.99 },
    }).as('createAlbum')

    cy.visit('/#/add')
  })

  it('displays the add album form', () => {
    cy.get('#title').should('be.visible')
    cy.get('#artist').should('be.visible')
    cy.get('#price').should('be.visible')
    cy.contains('button', 'Add Album').should('be.visible')
  })

  it('submits album and shows success message', () => {
    cy.get('#title').type('Blue Train')
    cy.get('#artist').type('John Coltrane')
    cy.get('#price').type('56.99')

    cy.contains('button', 'Add Album').click()

    cy.wait('@createAlbum').its('request.body').should('deep.equal', {
      title: 'Blue Train',
      artist: 'John Coltrane',
      price: 56.99,
    })

    cy.contains('Album added successfully')
    cy.get('#title').should('have.value', '')
    cy.get('#artist').should('have.value', '')
    cy.get('#price').should('have.value', '')
  })

  it('disables submit button while submitting', () => {
    cy.intercept('POST', '**/api/v1/album', (req) => {
      req.reply({ delay: 500, statusCode: 201, body: {} })
    }).as('createAlbum')

    cy.get('#title').type('Blue Train')
    cy.get('#artist').type('John Coltrane')
    cy.get('#price').type('56.99')

    cy.contains('button', 'Add Album').click()
    cy.contains('button', 'Adding...').should('be.visible')
    cy.contains('button', 'Add Album').should('be.visible')
  })

  it('shows error when API fails', () => {
    cy.intercept('POST', '**/api/v1/album', {
      statusCode: 400,
      body: { error: 'invalid album' },
    }).as('createAlbum')

    cy.get('#title').type('Test')
    cy.get('#artist').type('Artist')
    cy.get('#price').type('10')

    cy.contains('button', 'Add Album').click()

    cy.contains('Error: invalid album')
  })
})
