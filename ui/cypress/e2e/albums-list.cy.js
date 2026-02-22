const mockAlbums = [
  { id: 1, title: 'Blue Train', artist: 'John Coltrane', price: 56.99 },
  { id: 2, title: 'Giant Steps', artist: 'John Coltrane', price: 63.99 },
]

describe('Albums List', () => {
  beforeEach(() => {
    cy.intercept('GET', '**/api/v1/albums', {
      statusCode: 200,
      body: mockAlbums,
    }).as('getAlbums')

    cy.visit('/#/albums')
  })

  it('shows fetch button and prompt before loading', () => {
    cy.contains('button', 'Fetch Albums from REST API')
    cy.contains('Click the button to fetch albums from REST API')
  })

  it('fetches and displays albums when button clicked', () => {
    cy.contains('button', 'Fetch Albums from REST API').click()

    cy.wait('@getAlbums')

    cy.get('table').should('be.visible')
    cy.get('table th').contains('Id')
    cy.get('table th').contains('Title')
    cy.get('table th').contains('Artist')
    cy.get('table th').contains('Price')

    cy.get('table tbody tr').should('have.length', 2)
    cy.contains('Blue Train')
    cy.contains('John Coltrane')
    cy.contains('56.99')
  })

  it('shows loading state while fetching', () => {
    cy.intercept('GET', '**/api/v1/albums', (req) => {
      req.reply({ delay: 500, statusCode: 200, body: mockAlbums })
    }).as('getAlbums')

    cy.contains('button', 'Fetch Albums from REST API').click()

    cy.contains('button', 'Loading...').should('be.visible')
    cy.wait('@getAlbums')
    cy.contains('button', 'Fetch Albums from REST API').should('be.visible')
  })

  it('shows error when fetch fails', () => {
    cy.intercept('GET', '**/api/v1/albums', {
      statusCode: 500,
      body: {},
    }).as('getAlbums')

    cy.contains('button', 'Fetch Albums from REST API').click()

    cy.wait('@getAlbums')
    cy.contains('Error:')
  })
})
