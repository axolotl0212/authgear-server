import React, { useCallback, useContext, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { Dropdown, Label } from "@fluentui/react";
import deepEqual from "deep-equal";
import { Context, FormattedMessage } from "@oursky/react-messageformat";

import NavBreadcrumb from "../../NavBreadcrumb";
import {
  ModifiedIndicatorPortal,
  ModifiedIndicatorWrapper,
} from "../../ModifiedIndicatorPortal";
import ShowLoading from "../../ShowLoading";
import ShowError from "../../ShowError";
import FormTextField from "../../FormTextField";
import { passwordFieldErrorRules } from "../../PasswordField";
import AddIdentityForm from "./AddIdentityForm";
import {
  useDropdown,
  useIntegerTextField,
  useTextField,
} from "../../hook/useInput";
import { useAppConfigQuery } from "../portal/query/appConfigQuery";
import { UserQuery_node_User } from "./query/__generated__/UserQuery";
import { useUserQuery } from "./query/userQuery";
import { useCreateLoginIDIdentityMutation } from "./mutations/createIdentityMutation";
import { PortalAPIAppConfig } from "../../types";
import { useValidationError } from "../../error/useValidationError";
import { useGenericError } from "../../error/useGenericError";
import ShowUnhandledValidationErrorCause from "../../error/ShowUnhandledValidationErrorCauses";
import { FormContext } from "../../error/FormContext";
import { getActiveCountryCallingCode } from "../../util/countryCallingCode";

import styles from "./AddPhoneScreen.module.scss";

interface AddPhoneFormProps {
  appConfig: PortalAPIAppConfig | null;
  user: UserQuery_node_User | null;
}

const AddPhoneForm: React.FC<AddPhoneFormProps> = function AddPhoneForm(
  props: AddPhoneFormProps
) {
  const { appConfig, user } = props;
  const { userID } = useParams();
  const { renderToString } = useContext(Context);

  const {
    createIdentity,
    loading: creatingIdentity,
    error: createIdentityError,
  } = useCreateLoginIDIdentityMutation(userID);

  const countryCodeConfig = useMemo(() => {
    const countryCodeConfig = appConfig?.ui?.country_calling_code;
    const allowList = countryCodeConfig?.allowlist ?? [];
    const pinnedList = countryCodeConfig?.pinned_list ?? [];
    const values = getActiveCountryCallingCode(pinnedList, allowList);
    const defaultCallingCode = values[0];
    return {
      values,
      defaultCallingCode,
    };
  }, [appConfig]);

  const initialFormData = useMemo(() => {
    return {
      phone: "",
      countryCode: countryCodeConfig.defaultCallingCode,
      password: "",
    };
  }, [countryCodeConfig]);

  const [formData, setFormData] = useState(initialFormData);

  const [localValidationErrorMessage, setLocalViolationErrorMessage] = useState<
    string | undefined
  >(undefined);

  const { phone, countryCode, password } = formData;

  const { onChange: onPhoneChange } = useIntegerTextField((value) => {
    setFormData((prev) => ({
      ...prev,
      phone: value,
    }));
  });

  const { onChange: onPasswordChange } = useTextField((value) => {
    setFormData((prev) => ({ ...prev, password: value }));
  });

  const displayCountryCode = useCallback((countryCode: string) => {
    return `+ ${countryCode}`;
  }, []);

  const {
    options: countryCodeOptions,
    onChange: onCountryCodeChange,
  } = useDropdown(
    countryCodeConfig.values,
    (option) => {
      setFormData((prev) => ({
        ...prev,
        countryCode: option,
      }));
    },
    countryCodeConfig.defaultCallingCode,
    displayCountryCode
  );

  const isFormModified = useMemo(() => {
    return !deepEqual(initialFormData, formData);
  }, [formData, initialFormData]);

  const combinedPhone = useMemo(() => {
    return `+${countryCode}${phone}`;
  }, [countryCode, phone]);

  const resetForm = useCallback(() => {
    setFormData(initialFormData);
    setLocalViolationErrorMessage(undefined);
  }, [initialFormData]);

  const {
    unhandledCauses: rawUnhandledCauses,
    otherError,
    value: formContextValue,
  } = useValidationError(createIdentityError);

  const {
    errorMessageMap,
    unrecognizedError,
    unhandledCauses,
  } = useGenericError(otherError, rawUnhandledCauses, [
    {
      reason: "InvariantViolated",
      kind: "DuplicatedIdentity",
      errorMessageID: "AddPhoneScreen.error.duplicated-phone-number",
      field: "phone",
    },
    ...passwordFieldErrorRules,
  ]);

  return (
    <FormContext.Provider value={formContextValue}>
      {unrecognizedError && <ShowError error={unrecognizedError} />}
      <ShowUnhandledValidationErrorCause causes={unhandledCauses} />
      <ModifiedIndicatorPortal
        resetForm={resetForm}
        isModified={isFormModified}
      />
      <AddIdentityForm
        className={styles.form}
        appConfig={appConfig}
        user={user}
        password={password}
        onPasswordChange={onPasswordChange}
        passwordFieldErrorMessage={
          localValidationErrorMessage ?? errorMessageMap.password
        }
        loginIdKey="phone"
        loginId={combinedPhone}
        isFormModified={isFormModified}
        createIdentity={createIdentity}
        creatingIdentity={creatingIdentity}
        onLocalErrorMessageChange={setLocalViolationErrorMessage}
        loginIdField={
          <section className={styles.phoneNumberFields}>
            <Label className={styles.phoneNumberLabel}>
              <FormattedMessage id="AddPhoneScreen.phone.label" />
            </Label>
            <Dropdown
              className={styles.countryCode}
              options={countryCodeOptions}
              selectedKey={countryCode}
              onChange={onCountryCodeChange}
              ariaLabel={renderToString("AddPhoneScreen.country-code.label")}
            />
            <FormTextField
              jsonPointer=""
              parentJSONPointer=""
              fieldName="phone"
              fieldNameMessageID="AddPhoneScreen.phone.label"
              hideLabel={true}
              className={styles.phone}
              value={phone}
              onChange={onPhoneChange}
              ariaLabel={renderToString("AddPhoneScreen.phone.label")}
              errorMessage={errorMessageMap.phone}
            />
          </section>
        }
      />
    </FormContext.Provider>
  );
};

const AddPhoneScreen: React.FC = function AddPhoneScreen() {
  const { appID, userID } = useParams();

  const {
    user,
    loading: loadingUser,
    error: userError,
    refetch: refetchUser,
  } = useUserQuery(userID);
  const {
    effectiveAppConfig,
    loading: loadingAppConfig,
    error: appConfigError,
    refetch: refetchAppConfig,
  } = useAppConfigQuery(appID);

  const navBreadcrumbItems = useMemo(() => {
    return [
      { to: "../../..", label: <FormattedMessage id="UsersScreen.title" /> },
      { to: "../", label: <FormattedMessage id="UserDetailsScreen.title" /> },
      { to: ".", label: <FormattedMessage id="AddPhoneScreen.title" /> },
    ];
  }, []);

  if (loadingUser || loadingAppConfig) {
    return <ShowLoading />;
  }

  if (userError != null) {
    return <ShowError error={userError} onRetry={refetchUser} />;
  }

  if (appConfigError != null) {
    return <ShowError error={appConfigError} onRetry={refetchAppConfig} />;
  }

  return (
    <div className={styles.root}>
      <ModifiedIndicatorWrapper className={styles.content}>
        <NavBreadcrumb items={navBreadcrumbItems} />
        <AddPhoneForm appConfig={effectiveAppConfig} user={user} />
      </ModifiedIndicatorWrapper>
    </div>
  );
};

export default AddPhoneScreen;
