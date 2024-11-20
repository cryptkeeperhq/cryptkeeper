import React from 'react';

const InputField = ({ input, inputFieldTexts, handleInputChange, handleInputButtonClick }) => {
  const renderInput = () => {
    switch (input.inputType || 'textbox') {
      case 'textbox':
        return (
          <>
            <input
              type="text"
              className="form-control"
              name={input.id}
              value={inputFieldTexts[input.id] || ''}
              onChange={(e) => handleInputChange(e, input.id)}
            />
            <button
              style={{ position: "absolute", top: "-5px", right: "10px" }}
              className='btn btn-link'
              onClick={() => handleInputButtonClick(input.id)}
            >
              <i className='fa fa-list'></i>
            </button>
          </>
        );
      case 'select':
        return (
          <select
            className="form-select"
            name={input.id}
            value={inputFieldTexts[input.id] || ''}
            onChange={(e) => handleInputChange(e, input.id)}
          >
            {input.values.map((value) => (
              <option key={value} value={value}>
                {value}
              </option>
            ))}
          </select>
        );
      case 'textarea':
        return (
          <>
            <textarea
              className="form-control"
              name={input.id}
              value={inputFieldTexts[input.id] || ''}
              onChange={(e) => handleInputChange(e, input.id)}
            />
            <button
              style={{ position: "absolute", top: "-5px", right: "10px" }}
              className='btn btn-link'
              onClick={() => handleInputButtonClick(input.id)}
            >
              <i className='fa fa-list'></i>
            </button>
          </>
        );
      case 'checkbox':
        return (
          <div>
            {input.values.map((value) => (
              <div key={value} className="form-check">
                <input
                  type="checkbox"
                  className="form-check-input"
                  name={input.id}
                  value={value}
                  checked={inputFieldTexts[input.id]?.includes(value) || false}
                  onChange={(e) => handleInputChange(e, input.id)}
                />
                <label className="form-check-label">{value}</label>
              </div>
            ))}
          </div>
        );
      case 'radio':
        return (
          input.values.map((value) => (
            <div key={value} className="form-check form-check-inline">
              <input
                type="radio"
                className="form-check-input"
                name={input.id}
                value={value}
                checked={inputFieldTexts[input.id] === value}
                onChange={(e) => handleInputChange(e, input.id)}
              />
              <label className="form-check-label">{value}</label>
            </div>
          ))
        );
      default:
        return null;
    }
  };

  return (
    <div className="p-2 mb-1" style={{ position: "relative" }}>
      <div><b><label className="form-label">{input.label}</label></b></div>
      {renderInput()}
    </div>
  );
};

export default InputField;
