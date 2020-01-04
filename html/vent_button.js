'use strict';

const e = React.createElement;

class VentButton extends React.Component {
  constructor(props) {
    super(props);
    this.state = { on: false };
  }

  render() {
    if (this.state.on) {
      return 'Vent turned ON';
    }

    return e(
      'button',
      { onClick: () => this.setState({ on: true }) },
      'Switch'
    );
  }
}

const domContainer = document.querySelector('#vent_button_container');
ReactDOM.render(e(VentButton), domContainer);